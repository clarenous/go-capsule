package types

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

// Evidence
type Evidence struct {
	Digest      []byte
	Source      []byte
	ValidScript []byte
}

func (evid *Evidence) Hash(txid Hash, index uint64) (hash Hash) {
	var buf bytes.Buffer
	var b8 [8]byte
	binary.LittleEndian.PutUint64(b8[:], index)
	buf.Write(txid[:])
	buf.Write(b8[:])
	buf.Write(evid.Digest)
	buf.Write(evid.Source)
	buf.Write(evid.ValidScript)

	h := sha3.New256()
	h.Write(buf.Bytes())
	b := h.Sum(nil)
	r := bytes.NewReader(b)

	hash.ReadFrom(r)
	return hash
}
