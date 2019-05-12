package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clarenous/go-capsule/protocol/types/pb"
	"golang.org/x/crypto/sha3"
	"io"
)

const (
	HashSize       = 32
	HashStringSize = HashSize * 2
)

type Hash [HashSize]byte

// EmptyHash represents a 256-bit hash.
var EmptyHash = sha3.Sum256(nil)

// MarshalText satisfies the TextMarshaler interface.
// It returns the bytes of h encoded in hex,
// for formats that can't hold arbitrary binary data.
// It never returns an error.
func (h Hash) MarshalText() ([]byte, error) {
	v := make([]byte, 64)
	hex.Encode(v, h[:])
	return v, nil
}

// UnmarshalText satisfies the TextUnmarshaler interface.
// It decodes hex data from b into h.
func (h *Hash) UnmarshalText(v []byte) error {
	var b Hash
	if len(v) != 64 {
		return fmt.Errorf("bad length hash string %d", len(v))
	}
	_, err := hex.Decode(b[:], v)
	*h = b
	return err
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// If b is a JSON-encoded null, it copies the zero-value into h. Othwerwise, it
// decodes hex data from b into h.
func (h *Hash) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		*h = Hash{}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return h.UnmarshalText([]byte(s))
}

// Bytes returns the byte representation
func (h Hash) Bytes() []byte {
	var b32 [32]byte
	copy(b32[:], h[:])
	return b32[:]
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func NewHashFromString(src string) (Hash, error) {
	if len(src) > HashStringSize {
		return Hash{}, errors.New("larger than valid hash string size")
	}
	b, err := hex.DecodeString(src)
	if err != nil {
		return Hash{}, err
	}
	var b32 Hash
	copy(b32[:], b)
	return b32, nil
}

// WriteTo satisfies the io.WriterTo interface.
func (h Hash) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(h.Bytes())
	return int64(n), err
}

// ReadFrom satisfies the io.ReaderFrom interface.
func (h *Hash) ReadFrom(r io.Reader) (int64, error) {
	var b32 Hash
	n, err := io.ReadFull(r, b32[:])
	if err != nil {
		return int64(n), err
	}
	*h = b32
	return int64(n), nil
}

// IsZero tells whether a Hash pointer is nil or points to an all-zero
// hash.
func (h *Hash) IsZero() bool {
	if h == nil {
		return true
	}
	return *h == Hash{}
}

func (h Hash) Ptr() *Hash {
	return &h
}

func (h *Hash) Value() Hash {
	return *h
}

func (h *Hash) ToProto() *typespb.Hash {
	pb := typespb.NewHash(*h)
	return &pb
}

func (h *Hash) FromProto(pb *typespb.Hash) error {
	*h = pb.Byte32()
	return nil
}

func NewHashFromProto(pb *typespb.Hash) *Hash {
	h := new(Hash)
	h.FromProto(pb)
	return h
}
