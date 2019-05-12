package wallet

import (
	"bytes"
	"encoding/binary"
	"github.com/clarenous/go-capsule/crypto"
	"github.com/clarenous/go-capsule/crypto/ed25519"
	"github.com/clarenous/go-capsule/protocol/types"
)

const (
	RedeemScriptPrefixLength = 4
	UnlockScriptPrefixLength = 4
)

func createScriptHash(publicKey ed25519.PublicKey) (scriptHash types.Hash160, redeemScript []byte) {
	redeemScript = make([]byte, UnlockScriptPrefixLength+ed25519.PublicKeySize)
	binary.LittleEndian.PutUint32(redeemScript, ed25519.PublicKeySize)
	copy(redeemScript[UnlockScriptPrefixLength:], publicKey)
	h160 := crypto.Ripemd160(crypto.Sha256(redeemScript))
	(&scriptHash).SetBytes(h160)

	return scriptHash, redeemScript
}

func createUnlockScript(sig []byte) []byte {
	script := make([]byte, RedeemScriptPrefixLength+len(sig))
	binary.LittleEndian.PutUint32(script, uint32(len(sig)))
	copy(script[RedeemScriptPrefixLength:], sig)
	return script
}

func signTransaction(tx *types.Tx, sk ed25519.PrivateKey) []byte {
	var buf bytes.Buffer

	var b8 [8]byte
	var getUint64Bytes = func(n uint64) []byte {
		binary.LittleEndian.PutUint64(b8[:], tx.Version)
		return b8[:]
	}

	buf.Write(getUint64Bytes(tx.Version))

	for _, in := range tx.Inputs {
		buf.Write(in.ValueSource.Hash().Bytes())
		buf.Write(in.RedeemScript)
		// never write unlock script
		buf.Write(getUint64Bytes(in.Sequence))
	}

	for i, _ := range tx.Outputs {
		buf.Write(tx.OutHash(i).Bytes())
	}

	for _, evid := range tx.Evidences {
		buf.Write(evid.Digest)
		buf.Write(evid.Source)
		buf.Write(evid.ValidScript)
	}

	buf.Write(getUint64Bytes(tx.LockTime))

	return ed25519.Sign(sk, buf.Bytes())
}
