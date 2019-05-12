package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clarenous/go-capsule/protocol/types/pb"
	"io"
)

const (
	Hash160Size       = 20
	Hash160StringSize = Hash160Size * 2
)

type Hash160 [Hash160Size]byte

// MarshalText satisfies the TextMarshaler interface.
// It returns the bytes of h encoded in hex,
// for formats that can't hold arbitrary binary data.
// It never returns an error.
func (h Hash160) MarshalText() ([]byte, error) {
	v := make([]byte, Hash160StringSize)
	hex.Encode(v, h[:])
	return v, nil
}

// UnmarshalText satisfies the TextUnmarshaler interface.
// It decodes hex data from b into h.
func (h *Hash160) UnmarshalText(v []byte) error {
	var b Hash160
	if len(v) != Hash160StringSize {
		return fmt.Errorf("bad length Hash160 string %d", len(v))
	}
	_, err := hex.Decode(b[:], v)
	*h = b
	return err
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
// If b is a JSON-encoded null, it copies the zero-value into h. Othwerwise, it
// decodes hex data from b into h.
func (h *Hash160) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		*h = Hash160{}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return h.UnmarshalText([]byte(s))
}

// Bytes returns the byte representation
func (h Hash160) Bytes() []byte {
	return h[:]
}

func (h *Hash160) SetBytes(bs []byte) *Hash160 {
	var hp Hash160
	copy(hp[:], bs)

	*h = hp
	return h
}

func (h Hash160) String() string {
	return hex.EncodeToString(h[:])
}

func NewHash160FromString(src string) (Hash160, error) {
	if len(src) > Hash160StringSize {
		return Hash160{}, errors.New("larger than valid Hash160 string size")
	}
	b, err := hex.DecodeString(src)
	if err != nil {
		return Hash160{}, err
	}
	var b20 Hash160
	copy(b20[:], b)
	return b20, nil
}

// WriteTo satisfies the io.WriterTo interface.
func (h Hash160) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(h.Bytes())
	return int64(n), err
}

// ReadFrom satisfies the io.ReaderFrom interface.
func (h *Hash160) ReadFrom(r io.Reader) (int64, error) {
	var b20 Hash160
	n, err := io.ReadFull(r, b20[:])
	if err != nil {
		return int64(n), err
	}
	*h = b20
	return int64(n), nil
}

// IsZero tells whether a Hash160 pointer is nil or points to an all-zero
// hash.
func (h *Hash160) IsZero() bool {
	if h == nil {
		return true
	}
	return *h == Hash160{}
}

func (h Hash160) Ptr() *Hash160 {
	return &h
}

func (h *Hash160) Value() Hash160 {
	return *h
}

func (h *Hash160) ToProto() *typespb.Hash160 {
	pb := typespb.NewHash160(*h)
	return &pb
}

func (h *Hash160) FromProto(pb *typespb.Hash160) error {
	*h = pb.Byte20()
	return nil
}

func NewHash160FromProto(pb *typespb.Hash160) *Hash160 {
	h := new(Hash160)
	h.FromProto(pb)
	return h
}
