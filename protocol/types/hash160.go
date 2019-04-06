package types

import (
	"bytes"
	"encoding/binary"
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

type Hash160 typespb.Hash160

func to20Byte(b []byte) [20]byte {
	var b20 [20]byte
	copy(b20[:], b[:20])
	return b20
}

// NewHash160 convert the input byte array to hash
func NewHash160(b20 [20]byte) (h Hash160) {
	h.S0 = binary.BigEndian.Uint32(b20[0:4])
	h.S1 = binary.BigEndian.Uint32(b20[4:8])
	h.S2 = binary.BigEndian.Uint32(b20[8:12])
	h.S3 = binary.BigEndian.Uint32(b20[12:16])
	h.S4 = binary.BigEndian.Uint32(b20[16:20])
	return h
}

// Byte20 return the byte array representation
func (h Hash160) Byte20() (b20 [20]byte) {
	binary.BigEndian.PutUint32(b20[0:4], h.S0)
	binary.BigEndian.PutUint32(b20[4:8], h.S1)
	binary.BigEndian.PutUint32(b20[8:12], h.S2)
	binary.BigEndian.PutUint32(b20[12:16], h.S3)
	binary.BigEndian.PutUint32(b20[16:20], h.S4)
	return b20
}

// MarshalText satisfies the TextMarshaler interface.
// It returns the bytes of h encoded in hex,
// for formats that can't hold arbitrary binary data.
// It never returns an error.
func (h Hash160) MarshalText() ([]byte, error) {
	b := h.Byte20()
	v := make([]byte, 64)
	hex.Encode(v, b[:])
	return v, nil
}

// UnmarshalText satisfies the TextUnmarshaler interface.
// It decodes hex data from b into h.
func (h *Hash160) UnmarshalText(v []byte) error {
	var b [20]byte
	if len(v) != 40 {
		return fmt.Errorf("bad length Hash160 string %d", len(v))
	}
	_, err := hex.Decode(b[:], v)
	*h = NewHash160(b)
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
	b32 := h.Byte20()
	return b32[:]
}

func (h Hash160) String() string {
	b := h.Byte20()
	return hex.EncodeToString(b[:])
}

func NewHash160FromString(src string) (Hash160, error) {
	if len(src) > Hash160StringSize {
		return Hash160{}, errors.New("larger than valid Hash160 string size")
	}
	b, err := hex.DecodeString(src)
	if err != nil {
		return Hash160{}, err
	}
	var b20 [20]byte
	copy(b20[:], b)
	return NewHash160(b20), nil
}

// WriteTo satisfies the io.WriterTo interface.
func (h Hash160) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(h.Bytes())
	return int64(n), err
}

// ReadFrom satisfies the io.ReaderFrom interface.
func (h *Hash160) ReadFrom(r io.Reader) (int64, error) {
	var b20 [20]byte
	n, err := io.ReadFull(r, b20[:])
	if err != nil {
		return int64(n), err
	}
	*h = NewHash160(b20)
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

func (h *Hash160) Ptr() *Hash160 {
	return h
}
