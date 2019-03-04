package leveldb

import (
	"fmt"
	"sync"

	"github.com/golang/groupcache/lru"
	"github.com/golang/groupcache/singleflight"


	"github.com/clarenous/go-capsule/protocol/types"
)

const maxCachedBlocks = 30

func newBlockCache(fillFn func(hash *types.Hash) *types.Block) blockCache {
	return blockCache{
		lru:    lru.New(maxCachedBlocks),
		fillFn: fillFn,
	}
}

type blockCache struct {
	mu     sync.Mutex
	lru    *lru.Cache
	fillFn func(hash *types.Hash) *types.Block
	single singleflight.Group
}

func (c *blockCache) lookup(hash *types.Hash) (*types.Block, error) {
	if b, ok := c.get(hash); ok {
		return b, nil
	}

	block, err := c.single.Do(hash.String(), func() (interface{}, error) {
		b := c.fillFn(hash)
		if b == nil {
			return nil, fmt.Errorf("There are no block with given hash %s", hash.String())
		}

		c.add(b)
		return b, nil
	})
	if err != nil {
		return nil, err
	}
	return block.(*types.Block), nil
}

func (c *blockCache) get(hash *types.Hash) (*types.Block, bool) {
	c.mu.Lock()
	block, ok := c.lru.Get(*hash)
	c.mu.Unlock()
	if block == nil {
		return nil, ok
	}
	return block.(*types.Block), ok
}

func (c *blockCache) add(block *types.Block) {
	c.mu.Lock()
	c.lru.Add(block.Hash(), block)
	c.mu.Unlock()
}
