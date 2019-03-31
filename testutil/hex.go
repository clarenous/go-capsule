package testutil

import (
	"bytes"
	"io"
	"testing"

	"github.com/clarenous/go-capsule/protocol/types"
)

func MustDecodeHash(s string) (h types.Hash) {
	if err := h.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return h
}

func Serialize(t *testing.T, wt io.WriterTo) []byte {
	var b bytes.Buffer
	if _, err := wt.WriteTo(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
