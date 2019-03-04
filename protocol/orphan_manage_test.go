package protocol

import (
	"testing"
	"time"

	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/bytom/testutil"
)

var testBlocks = []*types.Block{
	&types.Block{BlockHeader: types.BlockHeader{
		PreviousBlockHash: types.Hash{V0: 1},
		Nonce:             0,
	}},
	&types.Block{BlockHeader: types.BlockHeader{
		PreviousBlockHash: types.Hash{V0: 1},
		Nonce:             1,
	}},
	&types.Block{BlockHeader: types.BlockHeader{
		PreviousBlockHash: types.Hash{V0: 2},
		Nonce:             3,
	}},
}

var blockHashes = []types.Hash{}

func init() {
	for _, block := range testBlocks {
		blockHashes = append(blockHashes, block.Hash())
	}
}

func TestOrphanManageAdd(t *testing.T) {
	cases := []struct {
		before    *OrphanManage
		after     *OrphanManage
		addOrphan *types.Block
	}{
		{
			before: &OrphanManage{
				orphan:      map[types.Hash]*orphanBlock{},
				prevOrphans: map[types.Hash][]*types.Hash{},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			addOrphan: testBlocks[0],
		},
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			addOrphan: testBlocks[0],
		},
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
					blockHashes[1]: &orphanBlock{testBlocks[1], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0], &blockHashes[1]},
				},
			},
			addOrphan: testBlocks[1],
		},
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
					blockHashes[2]: &orphanBlock{testBlocks[2], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
					types.Hash{V0: 2}: []*types.Hash{&blockHashes[2]},
				},
			},
			addOrphan: testBlocks[2],
		},
	}

	for i, c := range cases {
		c.before.Add(c.addOrphan)
		for _, orphan := range c.before.orphan {
			orphan.expiration = time.Time{}
		}
		if !testutil.DeepEqual(c.before, c.after) {
			t.Errorf("case %d: got %v want %v", i, c.before, c.after)
		}
	}
}

func TestOrphanManageDelete(t *testing.T) {
	cases := []struct {
		before *OrphanManage
		after  *OrphanManage
		remove *types.Hash
	}{
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			remove: &blockHashes[1],
		},
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan:      map[types.Hash]*orphanBlock{},
				prevOrphans: map[types.Hash][]*types.Hash{},
			},
			remove: &blockHashes[0],
		},
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
					blockHashes[1]: &orphanBlock{testBlocks[1], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0], &blockHashes[1]},
				},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{testBlocks[0], time.Time{}},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			remove: &blockHashes[1],
		},
	}

	for i, c := range cases {
		c.before.delete(c.remove)
		if !testutil.DeepEqual(c.before, c.after) {
			t.Errorf("case %d: got %v want %v", i, c.before, c.after)
		}
	}
}

func TestOrphanManageExpire(t *testing.T) {
	cases := []struct {
		before *OrphanManage
		after  *OrphanManage
	}{
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{
						testBlocks[0],
						time.Unix(1633479700, 0),
					},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan:      map[types.Hash]*orphanBlock{},
				prevOrphans: map[types.Hash][]*types.Hash{},
			},
		},
		{
			before: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{
						testBlocks[0],
						time.Unix(1633479702, 0),
					},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
			after: &OrphanManage{
				orphan: map[types.Hash]*orphanBlock{
					blockHashes[0]: &orphanBlock{
						testBlocks[0],
						time.Unix(1633479702, 0),
					},
				},
				prevOrphans: map[types.Hash][]*types.Hash{
					types.Hash{V0: 1}: []*types.Hash{&blockHashes[0]},
				},
			},
		},
	}

	for i, c := range cases {
		c.before.orphanExpire(time.Unix(1633479701, 0))
		if !testutil.DeepEqual(c.before, c.after) {
			t.Errorf("case %d: got %v want %v", i, c.before, c.after)
		}
	}
}

func TestOrphanManageNumLimit(t *testing.T) {
	cases := []struct{
		addOrphanBlockNum int
		expectOrphanBlockNum int
	}{
		{
			addOrphanBlockNum: 10,
			expectOrphanBlockNum: 10,
		},
		{
			addOrphanBlockNum: numOrphanBlockLimit,
			expectOrphanBlockNum: numOrphanBlockLimit,
		},
		{
			addOrphanBlockNum: numOrphanBlockLimit + 1,
			expectOrphanBlockNum: numOrphanBlockLimit,
		},
		{
			addOrphanBlockNum: numOrphanBlockLimit + 10,
			expectOrphanBlockNum: numOrphanBlockLimit,
		},
	}

	for i, c := range cases {
		orphanManage := &OrphanManage{
			orphan:      map[types.Hash]*orphanBlock{},
			prevOrphans: map[types.Hash][]*types.Hash{},
		}
		for num := 0; num < c.addOrphanBlockNum; num++ {
			orphanManage.Add(&types.Block{BlockHeader: types.BlockHeader{Height: uint64(num)}})
		}
		if (len(orphanManage.orphan) != c.expectOrphanBlockNum) {
			t.Errorf("case %d: got %d want %d", i, len(orphanManage.orphan), c.expectOrphanBlockNum)
		}
	}
}
