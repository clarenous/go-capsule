package wallet

import (
	"github.com/clarenous/go-capsule/common"
	"github.com/clarenous/go-capsule/crypto/ed25519"
	"github.com/clarenous/go-capsule/protocol/types"
)

type Wallet struct {
	privateKeys   map[string]ed25519.PrivateKey        // Public Key to Private key
	publicKeys    map[common.Address]ed25519.PublicKey // Address to Public Key
	redeemScripts map[types.Hash][]byte                // Output ID to RedeemScript

}
