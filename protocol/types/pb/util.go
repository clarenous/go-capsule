package typespb

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var errInvalidValue = errors.New("invalid value")

func WriteElement(w io.Writer, element interface{}) error {
	switch e := element.(type) {
	case uint32:
		var b4 [4]byte
		binary.LittleEndian.PutUint32(b4[:], e)
		if _, err := w.Write(b4[:]); err != nil {
			return err
		}
		return nil

	case uint64:
		var b8 [8]byte
		binary.LittleEndian.PutUint64(b8[:], e)
		if _, err := w.Write(b8[:]); err != nil {
			return err
		}
		return nil

	case []byte:
		if _, err := w.Write(e); err != nil {
			return err
		}
		return nil
	}

	switch v := reflect.ValueOf(element); v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return WriteElement(w, v.Elem().Interface())

	case reflect.Struct:
		typ := v.Type()
		for i := 0; i < typ.NumField(); i++ {
			c := v.Field(i)
			if !c.CanInterface() {
				return errInvalidValue
			}
			if err := WriteElement(w, c.Interface()); err != nil {
				t := v.Type()
				f := t.Field(i)
				return fmt.Errorf("%v, writing struct field %d (%s.%s) for element", err, i, t.Name(), f.Name)
			}
		}
		return nil
	}

	return errInvalidValue
}

func MockBytes32() [32]byte {
	var b32 [32]byte
	copy(b32[:], MockLenBytes(32))
	return b32
}

func MockBytes20() [20]byte {
	var b20 [20]byte
	copy(b20[:], MockLenBytes(20))
	return b20
}

func MockLenBytes(n int) []byte {
	bs := make([]byte, n)
	crand.Read(bs)
	return bs
}
