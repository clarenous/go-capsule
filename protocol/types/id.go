package types

import (
	"bytes"
	"golang.org/x/crypto/sha3"
	"reflect"
)

type Entry interface {
	bytesForID() []byte
}

func GetID(e Entry) (hash Hash) {
	if e == nil {
		return hash
	}

	// Nil pointer; not the same as nil interface above
	if v := reflect.ValueOf(e); v.Kind() == reflect.Ptr && v.IsNil() {
		return hash
	}

	h := sha3.New256()
	h.Write(e.bytesForID())
	b := h.Sum(nil)
	r := bytes.NewReader(b)

	hash.ReadFrom(r)
	return hash
}
