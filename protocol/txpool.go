package protocol

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/groupcache/lru"
	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/event"
	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/clarenous/go-capsule/protocol/state"
)

// msg type
const (
	MsgNewTx = iota
	MsgRemoveTx
	logModule = "protocol"
)

var (
	maxCachedErrTxs = 1000
	maxMsgChSize    = 1000
	maxNewTxNum     = 10000
	maxOrphanNum    = 2000

	orphanTTL                = 10 * time.Minute
	orphanExpireScanInterval = 3 * time.Minute

	// ErrTransactionNotExist is the pre-defined error message
	ErrTransactionNotExist = errors.New("transaction are not existed in the mempool")
	// ErrPoolIsFull indicates the pool is full
	ErrPoolIsFull = errors.New("transaction pool reach the max number")
	// ErrDustTx indicates transaction is dust tx
	ErrDustTx = errors.New("transaction is dust tx")
)

type TxMsgEvent struct{ TxMsg *TxPoolMsg }

// TxDesc store tx and related info for mining strategy
type TxDesc struct {
	Tx         *types.Tx `json:"transaction"`
	Added      time.Time `json:"-"`
	StatusFail bool      `json:"status_fail"`
	Height     uint64    `json:"-"`
	Weight     uint64    `json:"-"`
}

// TxPoolMsg is use for notify pool changes
type TxPoolMsg struct {
	*TxDesc
	MsgType int
}

type orphanTx struct {
	*TxDesc
	expiration time.Time
}

// TxPool is use for store the unconfirmed transaction
type TxPool struct {
	lastUpdated     int64
	mtx             sync.RWMutex
	store           Store
	pool            map[types.Hash]*TxDesc
	utxo            map[types.Hash]*types.Tx
	orphans         map[types.Hash]*orphanTx
	orphansByPrev   map[types.Hash]map[types.Hash]*orphanTx
	errCache        *lru.Cache
	eventDispatcher *event.Dispatcher
}

// NewTxPool init a new TxPool
func NewTxPool(store Store, dispatcher *event.Dispatcher) *TxPool {
	tp := &TxPool{
		lastUpdated:     time.Now().Unix(),
		store:           store,
		pool:            make(map[types.Hash]*TxDesc),
		utxo:            make(map[types.Hash]*types.Tx),
		orphans:         make(map[types.Hash]*orphanTx),
		orphansByPrev:   make(map[types.Hash]map[types.Hash]*orphanTx),
		errCache:        lru.New(maxCachedErrTxs),
		eventDispatcher: dispatcher,
	}
	go tp.orphanExpireWorker()
	return tp
}

// AddErrCache add a failed transaction record to lru cache
func (tp *TxPool) AddErrCache(txHash *types.Hash, err error) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	tp.errCache.Add(txHash, err)
}

// ExpireOrphan expire all the orphans that before the input time range
func (tp *TxPool) ExpireOrphan(now time.Time) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	for hash, orphan := range tp.orphans {
		if orphan.expiration.Before(now) {
			tp.removeOrphan(&hash)
		}
	}
}

// GetErrCache return the error of the transaction
func (tp *TxPool) GetErrCache(txHash *types.Hash) error {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	v, ok := tp.errCache.Get(txHash)
	if !ok {
		return nil
	}
	return v.(error)
}

// RemoveTransaction remove a transaction from the pool
func (tp *TxPool) RemoveTransaction(txHash *types.Hash) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	txD, ok := tp.pool[*txHash]
	if !ok {
		return
	}

	for i, _ := range txD.Tx.Outputs {
		delete(tp.utxo, txD.Tx.OutHash(i))
	}
	delete(tp.pool, *txHash)

	atomic.StoreInt64(&tp.lastUpdated, time.Now().Unix())
	tp.eventDispatcher.Post(TxMsgEvent{TxMsg: &TxPoolMsg{TxDesc: txD, MsgType: MsgRemoveTx}})
	log.WithFields(log.Fields{"module": logModule, "tx_id": txHash}).Debug("remove tx from mempool")
}

// GetTransaction return the TxDesc by hash
func (tp *TxPool) GetTransaction(txHash *types.Hash) (*TxDesc, error) {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	if txD, ok := tp.pool[*txHash]; ok {
		return txD, nil
	}
	return nil, ErrTransactionNotExist
}

// GetTransactions return all the transactions in the pool
func (tp *TxPool) GetTransactions() []*TxDesc {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	txDs := make([]*TxDesc, len(tp.pool))
	i := 0
	for _, desc := range tp.pool {
		txDs[i] = desc
		i++
	}
	return txDs
}

// IsTransactionInPool check wheather a transaction in pool or not
func (tp *TxPool) IsTransactionInPool(txHash *types.Hash) bool {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	_, ok := tp.pool[*txHash]
	return ok
}

// IsTransactionInErrCache check wheather a transaction in errCache or not
func (tp *TxPool) IsTransactionInErrCache(txHash *types.Hash) bool {
	tp.mtx.RLock()
	defer tp.mtx.RUnlock()

	_, ok := tp.errCache.Get(txHash)
	return ok
}

// HaveTransaction IsTransactionInErrCache check is  transaction in errCache or pool
func (tp *TxPool) HaveTransaction(txHash *types.Hash) bool {
	return tp.IsTransactionInPool(txHash) || tp.IsTransactionInErrCache(txHash)
}

func (tp *TxPool) IsDust(tx *types.Tx) bool {
	// TODO: 增加粉尘交易规则
	return false
}

func (tp *TxPool) processTransaction(tx *types.Tx, height uint64) (bool, error) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()

	txD := &TxDesc{
		Tx:     tx,
		Weight: tx.SerializedSize(),
		Height: height,
	}
	requireParents, err := tp.checkOrphanUtxos(tx)
	if err != nil {
		return false, err
	}

	if len(requireParents) > 0 {
		return true, tp.addOrphan(txD, requireParents)
	}

	if err := tp.addTransaction(txD); err != nil {
		return false, err
	}

	tp.processOrphans(txD)
	return false, nil
}

// ProcessTransaction is the main entry for txpool handle new tx, ignore dust tx.
func (tp *TxPool) ProcessTransaction(tx *types.Tx, height uint64) (bool, error) {
	if tp.IsDust(tx) {
		log.WithFields(log.Fields{"module": logModule, "tx_id": tx.Hash().String()}).Warn("dust tx")
		return false, nil
	}
	return tp.processTransaction(tx, height)
}

func (tp *TxPool) addOrphan(txD *TxDesc, requireParents []*types.Hash) error {
	if len(tp.orphans) >= maxOrphanNum {
		return ErrPoolIsFull
	}

	orphan := &orphanTx{txD, time.Now().Add(orphanTTL)}
	tp.orphans[txD.Tx.Hash()] = orphan
	for _, hash := range requireParents {
		if _, ok := tp.orphansByPrev[*hash]; !ok {
			tp.orphansByPrev[*hash] = make(map[types.Hash]*orphanTx)
		}
		tp.orphansByPrev[*hash][txD.Tx.Hash()] = orphan
	}
	return nil
}

func (tp *TxPool) addTransaction(txD *TxDesc) error {
	if len(tp.pool) >= maxNewTxNum {
		return ErrPoolIsFull
	}

	tx := txD.Tx
	txD.Added = time.Now()
	tp.pool[tx.Hash()] = txD
	for i, _ := range tx.Outputs {
		tp.utxo[tx.OutHash(i)] = tx
	}

	atomic.StoreInt64(&tp.lastUpdated, time.Now().Unix())
	tp.eventDispatcher.Post(TxMsgEvent{TxMsg: &TxPoolMsg{TxDesc: txD, MsgType: MsgNewTx}})
	log.WithFields(log.Fields{"module": logModule, "tx_id": tx.Hash().String()}).Debug("Add tx to mempool")
	return nil
}

func (tp *TxPool) checkOrphanUtxos(tx *types.Tx) ([]*types.Hash, error) {
	view := state.NewUtxoViewpoint()
	if err := tp.store.GetTransactionsUtxo(view, []*types.Tx{tx}); err != nil {
		return nil, err
	}

	parents := []*types.Hash{}
	for _, in := range tx.Inputs {
		hash := in.ValueSource.Hash()
		if !view.CanSpend(&hash) && tp.utxo[hash] == nil {
			parents = append(parents, &in.ValueSource.TxID)
		}
	}
	return parents, nil
}

func (tp *TxPool) orphanExpireWorker() {
	ticker := time.NewTicker(orphanExpireScanInterval)
	for now := range ticker.C {
		tp.ExpireOrphan(now)
	}
}

func (tp *TxPool) processOrphans(txD *TxDesc) {
	processOrphans := []*orphanTx{}
	addRely := func(tx *types.Tx) {
		parentHash := tx.Hash()
		if orphans, ok := tp.orphansByPrev[parentHash]; ok {
			for _, orphan := range orphans {
				processOrphans = append(processOrphans, orphan)
			}
			delete(tp.orphansByPrev, parentHash)
		}
	}

	addRely(txD.Tx)
	for ; len(processOrphans) > 0; processOrphans = processOrphans[1:] {
		processOrphan := processOrphans[0]
		requireParents, err := tp.checkOrphanUtxos(processOrphan.Tx)
		if err != nil {
			log.WithFields(log.Fields{"module": logModule, "err": err}).Error("processOrphans got unexpected error")
			continue
		}

		if len(requireParents) == 0 {
			addRely(processOrphan.Tx)
			tp.removeOrphan(processOrphan.Tx.Hash().Ptr())
			tp.addTransaction(processOrphan.TxDesc)
		}
	}
}

func (tp *TxPool) removeOrphan(hash *types.Hash) {
	orphan, ok := tp.orphans[*hash]
	if !ok {
		return
	}

	for _, in := range orphan.Tx.Inputs {
		parentHash := in.ValueSource.TxID
		orphans, ok := tp.orphansByPrev[parentHash]
		if !ok {
			continue
		}

		if delete(orphans, *hash); len(orphans) == 0 {
			delete(tp.orphansByPrev, parentHash)
		}
	}
	delete(tp.orphans, *hash)
}
