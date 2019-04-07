package miningpool

import (
	"errors"
	"github.com/clarenous/go-capsule/event"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/mining"
	"github.com/clarenous/go-capsule/protocol"
	"github.com/clarenous/go-capsule/protocol/types"
)

const (
	maxSubmitChSize = 50
)

type submitBlockMsg struct {
	blockHeader *types.BlockHeader
	reply       chan error
}

// MiningPool is the support struct for p2p mine pool
type MiningPool struct {
	mutex    sync.RWMutex
	block    *types.Block
	submitCh chan *submitBlockMsg

	chain           *protocol.Chain
	txPool          *protocol.TxPool
	eventDispatcher *event.Dispatcher
}

// NewMiningPool will create a new MiningPool
func NewMiningPool(c *protocol.Chain, txPool *protocol.TxPool, dispatcher *event.Dispatcher) *MiningPool {
	m := &MiningPool{
		submitCh:        make(chan *submitBlockMsg, maxSubmitChSize),
		chain:           c,
		txPool:          txPool,
		eventDispatcher: dispatcher,
	}
	m.generateBlock()
	go m.blockUpdater()
	return m
}

// blockUpdater is the goroutine for keep update mining block
func (m *MiningPool) blockUpdater() {
	for {
		select {
		case <-m.chain.BlockWaiter(m.chain.BestBlockHeight() + 1):
			m.generateBlock()

		case submitMsg := <-m.submitCh:
			err := m.submitWork(submitMsg.blockHeader)
			if err == nil {
				m.generateBlock()
			}
			submitMsg.reply <- err
		}
	}
}

// generateBlock generates a block template to mine
func (m *MiningPool) generateBlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	block, err := mining.NewBlockTemplate(m.chain, m.txPool)
	if err != nil {
		log.Errorf("miningpool: failed on create NewBlockTemplate: %v", err)
		return
	}
	m.block = block
}

// GetWork will return a block header for p2p mining
func (m *MiningPool) GetWork() (*types.BlockHeader, error) {
	if m.block != nil {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
		bh := m.block.BlockHeader
		return &bh, nil
	}
	return nil, errors.New("no block is ready for mining")
}

// SubmitWork will try to submit the result to the blockchain
func (m *MiningPool) SubmitWork(bh *types.BlockHeader) error {
	reply := make(chan error, 1)
	m.submitCh <- &submitBlockMsg{blockHeader: bh, reply: reply}
	err := <-reply
	if err != nil {
		log.WithFields(log.Fields{"err": err, "height": bh.Height}).Warning("submitWork failed")
	}
	return err
}

func (m *MiningPool) submitWork(bh *types.BlockHeader) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.block == nil || bh.Previous != m.block.Previous {
		return errors.New("pending mining block has been changed")
	}

	m.block.Proof = bh.Proof
	m.block.Timestamp = bh.Timestamp
	isOrphan, err := m.chain.ProcessBlock(m.block)
	if err != nil {
		return err
	}
	if isOrphan {
		return errors.New("submit result is orphan")
	}

	if err := m.eventDispatcher.Post(event.NewMinedBlockEvent{Block: m.block}); err != nil {
		return err
	}

	return nil
}
