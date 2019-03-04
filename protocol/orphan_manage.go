package protocol

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/protocol/types"
)

var (
	orphanBlockTTL           = 60 * time.Minute
	orphanExpireWorkInterval = 3 * time.Minute
	numOrphanBlockLimit      = 256
)

type orphanBlock struct {
	*types.Block
	expiration time.Time
}

// OrphanManage is use to handle all the orphan block
type OrphanManage struct {
	orphan      map[types.Hash]*orphanBlock
	prevOrphans map[types.Hash][]*types.Hash
	mtx         sync.RWMutex
}

// NewOrphanManage return a new orphan block
func NewOrphanManage() *OrphanManage {
	o := &OrphanManage{
		orphan:      make(map[types.Hash]*orphanBlock),
		prevOrphans: make(map[types.Hash][]*types.Hash),
	}

	go o.orphanExpireWorker()
	return o
}

// BlockExist check is the block in OrphanManage
func (o *OrphanManage) BlockExist(hash *types.Hash) bool {
	o.mtx.RLock()
	_, ok := o.orphan[*hash]
	o.mtx.RUnlock()
	return ok
}

// Add will add the block to OrphanManage
func (o *OrphanManage) Add(block *types.Block) {
	blockHash := block.Hash()
	o.mtx.Lock()
	defer o.mtx.Unlock()

	if _, ok := o.orphan[blockHash]; ok {
		return
	}

	if len(o.orphan) >= numOrphanBlockLimit {
		log.WithFields(log.Fields{"module": logModule, "hash": blockHash.String(), "height": block.Height}).Info("the number of orphan blocks exceeds the limit")
		return
	}

	o.orphan[blockHash] = &orphanBlock{block, time.Now().Add(orphanBlockTTL)}
	o.prevOrphans[block.PreviousBlockHash] = append(o.prevOrphans[block.PreviousBlockHash], &blockHash)

	log.WithFields(log.Fields{"module": logModule, "hash": blockHash.String(), "height": block.Height}).Info("add block to orphan")
}

// Delete will delete the block from OrphanManage
func (o *OrphanManage) Delete(hash *types.Hash) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.delete(hash)
}

// Get return the orphan block by hash
func (o *OrphanManage) Get(hash *types.Hash) (*types.Block, bool) {
	o.mtx.RLock()
	block, ok := o.orphan[*hash]
	o.mtx.RUnlock()
	return block.Block, ok
}

// GetPrevOrphans return the list of child orphans
func (o *OrphanManage) GetPrevOrphans(hash *types.Hash) ([]*types.Hash, bool) {
	o.mtx.RLock()
	prevOrphans, ok := o.prevOrphans[*hash]
	o.mtx.RUnlock()
	return prevOrphans, ok
}

func (o *OrphanManage) delete(hash *types.Hash) {
	block, ok := o.orphan[*hash]
	if !ok {
		return
	}
	delete(o.orphan, *hash)

	prevOrphans, ok := o.prevOrphans[block.Block.PreviousBlockHash]
	if !ok || len(prevOrphans) == 1 {
		delete(o.prevOrphans, block.Block.PreviousBlockHash)
		return
	}

	for i, preOrphan := range prevOrphans {
		if preOrphan == hash {
			o.prevOrphans[block.Block.PreviousBlockHash] = append(prevOrphans[:i], prevOrphans[i+1:]...)
			return
		}
	}
}

func (o *OrphanManage) orphanExpireWorker() {
	ticker := time.NewTicker(orphanExpireWorkInterval)
	for now := range ticker.C {
		o.orphanExpire(now)
	}
	ticker.Stop()
}

func (o *OrphanManage) orphanExpire(now time.Time) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	for hash, orphan := range o.orphan {
		if orphan.expiration.Before(now) {
			o.delete(&hash)
		}
	}
}
