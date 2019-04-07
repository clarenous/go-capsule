package typespb

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

func TestWriteElement(t *testing.T) {
	bh := &BlockHeader{
		ChainId:         NewHash(MockBytes32()).Ptr(),
		Version:         1,
		Height:          100,
		Timestamp:       uint64(time.Now().Unix()),
		Previous:        NewHash(MockBytes32()).Ptr(),
		TransactionRoot: NewHash(MockBytes32()).Ptr(),
		WitnessRoot:     NewHash(MockBytes32()).Ptr(),
		Proof:           MockLenBytes(20),
	}

	var buf bytes.Buffer
	if err := WriteElement(&buf, bh); err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(hex.EncodeToString(buf.Bytes()))
}
