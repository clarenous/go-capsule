package types

import (
	crand "crypto/rand"
	ca "github.com/clarenous/go-capsule/consensus/algorithm"
	"math/rand"
	"time"
)

func MockBlock() *Block {
	var txsCount = rand.Intn(20)
	txs := make([]*Tx, txsCount)
	for i := range txs {
		txs[i] = MockTx()
	}
	header := MockBlockHeader()

	return &Block{
		BlockHeader:  *header,
		Transactions: txs,
	}
}

func MockBlockHeader() *BlockHeader {
	return &BlockHeader{
		ChainID:         MockHash(),
		Version:         1,
		Height:          rand.Uint64() % 1000000,
		Timestamp:       uint64(time.Now().Unix()) + rand.Uint64()%(3600*24*365),
		Previous:        MockHash(),
		TransactionRoot: MockHash(),
		WitnessRoot:     MockHash(),
		Proof:           MockProof(),
	}
}

func MockProof() ca.Proof {
	return nil
}

func MockTx() *Tx {
	var inputsCount, outputsCount, evidencesCount = rand.Intn(10), rand.Intn(5), rand.Intn(5)
	inputs, outputs, evidences := make([]TxIn, inputsCount), make([]TxOut, outputsCount), make([]Evidence, evidencesCount)
	for i := range inputs {
		in := MockTxIn()
		inputs[i] = *in
	}
	for i := range outputs {
		out := MockTxOut()
		outputs[i] = *out
	}
	for i := range evidences {
		evid := MockEvidence()
		evidences[i] = *evid
	}

	return &Tx{
		Version:   1,
		Inputs:    inputs,
		Outputs:   outputs,
		Evidences: evidences,
		LockTime:  0,
	}
}

func MockTxIn() *TxIn {
	return &TxIn{
		ValueSource: ValueSource{
			TxID:  MockHash(),
			Index: rand.Uint64() % 20,
		},
		RedeemScript: MockLenBytes(40),
		UnlockScript: MockLenBytes(80),
		Sequence:     1<<64 - 1,
	}
}

func MockTxOut() *TxOut {
	return &TxOut{
		Value:      rand.Uint64(),
		ScriptHash: MockHash160(),
	}
}

func MockEvidence() *Evidence {
	return &Evidence{
		Digest:      MockLenBytes(32),
		Source:      MockLenBytes(32),
		ValidScript: MockLenBytes(40),
	}
}

func MockHash() Hash {
	return Hash{
		S0: rand.Uint64(),
		S1: rand.Uint64(),
		S2: rand.Uint64(),
		S3: rand.Uint64(),
	}
}

func MockHash160() Hash160 {
	return Hash160{
		S0: rand.Uint32(),
		S1: rand.Uint32(),
		S2: rand.Uint32(),
		S3: rand.Uint32(),
		S4: rand.Uint32(),
	}
}

func MockLenBytes(n int) []byte {
	bs := make([]byte, n)
	crand.Read(bs)
	return bs
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
