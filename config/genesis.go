package config

import (
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
	txStatus := types.NewTransactionStatus()
	if err := txStatus.SetStatus(0, false); err != nil {
		log.Panicf(err.Error())
	}
	txStatusHash, err := types.TxStatusMerkleRoot(txStatus.VerifyStatus)
	if err != nil {
		log.Panicf("fail on calc genesis tx status merkle root")
	}

	merkleRoot, err := types.TxMerkleRoot([]*types.Tx{tx.Tx})
	if err != nil {
		log.Panicf("fail on calc genesis tx merkel root")
	}

	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Version:   1,
			Height:    0,
			Nonce:     9253507043297,
			Timestamp: 1524549600,
			Bits:      2161727821137910632,
			BlockCommitment: types.BlockCommitment{
				TransactionsMerkleRoot: merkleRoot,
				TransactionStatusHash:  txStatusHash,
			},
		},
		Transactions: []*types.Tx{tx},
	}
	return block
}

func testNetGenesisBlock() *types.Block {
	tx := genesisTx()
	txStatus := types.NewTransactionStatus()
	if err := txStatus.SetStatus(0, false); err != nil {
		log.Panicf(err.Error())
	}
	txStatusHash, err := types.TxStatusMerkleRoot(txStatus.VerifyStatus)
	if err != nil {
		log.Panicf("fail on calc genesis tx status merkle root")
	}

	merkleRoot, err := types.TxMerkleRoot([]*types.Tx{tx.Tx})
	if err != nil {
		log.Panicf("fail on calc genesis tx merkel root")
	}

	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Version:   1,
			Height:    0,
			Nonce:     9253507043297,
			Timestamp: 1528945000,
			Bits:      2305843009214532812,
			BlockCommitment: types.BlockCommitment{
				TransactionsMerkleRoot: merkleRoot,
				TransactionStatusHash:  txStatusHash,
			},
		},
		Transactions: []*types.Tx{tx},
	}
	return block
}

func soloNetGenesisBlock() *types.Block {
	tx := genesisTx()
	txStatus := types.NewTransactionStatus()
	if err := txStatus.SetStatus(0, false); err != nil {
		log.Panicf(err.Error())
	}
	txStatusHash, err := types.TxStatusMerkleRoot(txStatus.VerifyStatus)
	if err != nil {
		log.Panicf("fail on calc genesis tx status merkle root")
	}

	merkleRoot, err := types.TxMerkleRoot([]*types.Tx{tx.Tx})
	if err != nil {
		log.Panicf("fail on calc genesis tx merkel root")
	}

	block := &types.Block{
		BlockHeader: types.BlockHeader{
			Version:   1,
			Height:    0,
			Nonce:     9253507043297,
			Timestamp: 1528945000,
			Bits:      2305843009214532812,
			BlockCommitment: types.BlockCommitment{
				TransactionsMerkleRoot: merkleRoot,
				TransactionStatusHash:  txStatusHash,
			},
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
