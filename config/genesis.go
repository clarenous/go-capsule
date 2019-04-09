package config

import (
	ca "github.com/clarenous/go-capsule/consensus/algorithm"
	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/consensus"

	"github.com/clarenous/go-capsule/protocol/types"
)

func genesisTx() *types.Tx {
	tx := &types.Tx{
		Version: 1,
		Inputs:  []types.TxIn{
			//types.NewCoinbaseInput([]byte("Information is power. -- Jan/11/2013. Computing is power. -- Apr/24/2018.")),
		},
		Outputs: []types.TxOut{
			//types.NewTxOutput(*consensus.BTMAssetID, consensus.InitialBlockSubsidy, contract),
		},
	}
	return tx
}

func mainNetGenesisBlock() *types.Block {
	tx := genesisTx()

	merkleRoot, err := types.TxMerkleRoot([]*types.Tx{tx})
	if err != nil {
		log.Panicf("fail on calc genesis tx merkle root")
	}

	proof, err := ca.NewProof(consensus.ProofType, 2161727821137910632, 9253507043297)
	if err != nil {
		log.Panicf("fail on calc genesis proof")
	}

	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Version:         1,
			Height:          0,
			Timestamp:       1524549600,
			TransactionRoot: merkleRoot,
			WitnessRoot:     merkleRoot,
			Proof:           proof,
		},
		Transactions: []*types.Tx{tx},
	}
	return block
}

func testNetGenesisBlock() *types.Block {
	tx := genesisTx()

	merkleRoot, err := types.TxMerkleRoot([]*types.Tx{tx})
	if err != nil {
		log.Panicf("fail on calc genesis tx merkle root")
	}

	proof, _ := ca.NewProof(consensus.ProofType, 2305843009214532812, 9253507043297)

	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Version:         1,
			Height:          0,
			Timestamp:       1528945000,
			TransactionRoot: merkleRoot,
			WitnessRoot:     merkleRoot,
			Proof:           proof,
		},
		Transactions: []*types.Tx{tx},
	}
	return block
}

func soloNetGenesisBlock() *types.Block {
	tx := genesisTx()

	merkleRoot, err := types.TxMerkleRoot([]*types.Tx{tx})
	if err != nil {
		log.Panicf("fail on calc genesis tx merkle root")
	}

	proof, _ := ca.NewProof(consensus.ProofType, 2305843009214532812, 9253507043297)

	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Version:         1,
			Height:          0,
			Timestamp:       1528945000,
			TransactionRoot: merkleRoot,
			WitnessRoot:     merkleRoot,
			Proof:           proof,
		},
		Transactions: []*types.Tx{tx},
	}
	return block
}

// GenesisBlock will return genesis block
func GenesisBlock() *types.Block {
	return map[string]func() *types.Block{
		"main": mainNetGenesisBlock,
		"test": testNetGenesisBlock,
		"solo": soloNetGenesisBlock,
	}[consensus.ActiveNetParams.Name]()
}
