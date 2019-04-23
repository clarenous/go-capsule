package leveldb

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/clarenous/go-capsule/protocol/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	txLocPrefix   = []byte("TL:")
	evidLocPrefix = []byte("EVIDL:")
)

func (s *Store) GetTransaction(hash *types.Hash) (*types.Tx, error) {
	var prefix [35]byte
	copy(prefix[:], txLocPrefix)
	copy(prefix[3:], hash.Bytes())

	iter := dbm.IteratePrefix(s.db, prefix[:])
	for iter.Valid() {
		key := iter.Key()
		loc := types.NewTxLocFromBytes(key[3:])
		blkBytes := s.db.Get(calcBlockKey(&loc.BlockHash))
		if blkBytes != nil && len(blkBytes) >= int(loc.Offset+loc.Length) {
			txBytes := blkBytes[loc.Offset : loc.Offset+loc.Length]
			tx := new(types.Tx)
			if err := tx.UnmarshalText(txBytes); err == nil {
				return tx, nil
			}
		}
		iter.Next()
	}

	return nil, fmt.Errorf("fail to find transaction by hash %s", hash.String())
}

func (s *Store) saveTxLocs(batch dbm.Batch, locs []*types.TxLoc) {
	for _, loc := range locs {
		s.saveTxLoc(batch, loc)
	}
}

// saveTxLoc saves txLoc into batch
func (s *Store) saveTxLoc(batch dbm.Batch, loc *types.TxLoc) {
	var b83 [83]byte
	var b80 = loc.Byte80()

	copy(b83[:], txLocPrefix)
	copy(b83[3:], b80[:])

	batch.Set(b83[:], []byte{})
}

func (s *Store) GetEvidence(hash *types.Hash) (*types.Evidence, *types.Tx, int, error) {
	var prefix [38]byte
	var getIndex = func(b8 []byte) int {
		return int(binary.LittleEndian.Uint64(b8))
	}
	copy(prefix[:], evidLocPrefix)
	copy(prefix[6:], hash.Bytes())

	iter := dbm.IteratePrefix(s.db, prefix[:])
	for iter.Valid() {
		if key := iter.Key(); len(key) == 78 {
			fmt.Println("evidence key", key, hex.EncodeToString(key), string(key))
			var txHash types.Hash
			copy(txHash[:], key[38:70])
			if tx, err := s.GetTransaction(&txHash); err == nil {
				if index := getIndex(key[70:78]); len(tx.Evidences) > index {
					return &tx.Evidences[index], tx, index, nil
				}
			}
		}
		iter.Next()
	}

	return nil, nil, 0, fmt.Errorf("fail to find evidence by hash %s", hash.String())
}

// saveEvidLoc saves evidence loc into batch
func (s *Store) saveEvidLocs(batch dbm.Batch, blk *types.Block) {
	var b78 [78]byte
	var b8 [8]byte
	var getIndex = func(index int) []byte {
		binary.LittleEndian.PutUint64(b8[:], uint64(index))
		return b8[:]
	}
	copy(b78[:], evidLocPrefix)

	for _, tx := range blk.Transactions {
		txHash := tx.Hash()
		for j, evid := range tx.Evidences {
			var key [78]byte
			evidHash := evid.Hash(txHash, uint64(j))
			copy(b78[6:], evidHash[:])
			copy(b78[38:], txHash[:])
			copy(b78[70:], getIndex(j))
			copy(key[:], b78[:])
			batch.Set(key[:], []byte{})
		}
	}
}
