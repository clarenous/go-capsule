package netsync

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tendermint/go-amino"

	"github.com/clarenous/go-capsule/protocol/types"
)

func init() {
	RegisterAmino(cdc)
}

//protocol msg byte
const (
	BlockchainChannel = byte(0x40)

	BlockRequestByte    = byte(0x10)
	BlockResponseByte   = byte(0x11)
	HeadersRequestByte  = byte(0x12)
	HeadersResponseByte = byte(0x13)
	BlocksRequestByte   = byte(0x14)
	BlocksResponseByte  = byte(0x15)
	StatusRequestByte   = byte(0x20)
	StatusResponseByte  = byte(0x21)
	NewTransactionByte  = byte(0x30)
	NewMineBlockByte    = byte(0x40)
	FilterLoadByte      = byte(0x50)
	FilterAddByte       = byte(0x51)
	FilterClearByte     = byte(0x52)

	maxBlockchainResponseSize = 22020096 + 2
)

var cdc = amino.NewCodec()

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*BlockchainMessage)(nil), nil)
	cdc.RegisterConcrete(&GetBlockMessage{}, string(BlockRequestByte), nil)
	cdc.RegisterConcrete(&BlockMessage{}, string(BlockResponseByte), nil)
	cdc.RegisterConcrete(&GetHeadersMessage{}, string(HeadersRequestByte), nil)
	cdc.RegisterConcrete(&HeadersMessage{}, string(HeadersResponseByte), nil)
	cdc.RegisterConcrete(&GetBlocksMessage{}, string(BlocksRequestByte), nil)
	cdc.RegisterConcrete(&BlocksMessage{}, string(BlocksResponseByte), nil)
	cdc.RegisterConcrete(&StatusRequestMessage{}, string(StatusRequestByte), nil)
	cdc.RegisterConcrete(&StatusResponseMessage{}, string(StatusResponseByte), nil)
	cdc.RegisterConcrete(&TransactionMessage{}, string(NewTransactionByte), nil)
	cdc.RegisterConcrete(&MineBlockMessage{}, string(NewMineBlockByte), nil)
	cdc.RegisterConcrete(&FilterLoadMessage{}, string(FilterLoadByte), nil)
	cdc.RegisterConcrete(&FilterAddMessage{}, string(FilterAddByte), nil)
	cdc.RegisterConcrete(&FilterClearMessage{}, string(FilterClearByte), nil)
}

//BlockchainMessage is a generic message for this reactor.
type BlockchainMessage interface {
	String() string
}

//DecodeMessage decode msg
func DecodeMessage(bz []byte) (msgType byte, msg BlockchainMessage, err error) {
	msgType = bz[0]
	r := bytes.NewReader(bz)
	n := int64(0)
	n, err = cdc.UnmarshalBinaryLengthPrefixedReader(r, &msg, maxBlockchainResponseSize)
	if err != nil && int(n) != len(bz) {
		err = errors.New("DecodeMessage() had bytes left over")
	}
	return
}

//GetBlockMessage request blocks from remote peers by height/hash
type GetBlockMessage struct {
	Height  uint64
	RawHash [32]byte
}

//GetHash reutrn the hash of the request
func (m *GetBlockMessage) GetHash() *types.Hash {
	hash := types.Hash(m.RawHash)
	return &hash
}

func (m *GetBlockMessage) String() string {
	if m.Height > 0 {
		return fmt.Sprintf("{height: %d}", m.Height)
	}
	return fmt.Sprintf("{hash: %s}", hex.EncodeToString(m.RawHash[:]))
}

//BlockMessage response get block msg
type BlockMessage struct {
	RawBlock []byte
}

//NewBlockMessage construct bock response msg
func NewBlockMessage(block *types.Block) (*BlockMessage, error) {
	rawBlock, err := block.MarshalText()
	if err != nil {
		return nil, err
	}
	return &BlockMessage{RawBlock: rawBlock}, nil
}

//GetBlock get block from msg
func (m *BlockMessage) GetBlock() (*types.Block, error) {
	block := &types.Block{
		BlockHeader:  types.BlockHeader{},
		Transactions: []*types.Tx{},
	}
	if err := block.UnmarshalText(m.RawBlock); err != nil {
		return nil, err
	}
	return block, nil
}

func (m *BlockMessage) String() string {
	block, err := m.GetBlock()
	if err != nil {
		return "{err: wrong message}"
	}
	blockHash := block.Hash()
	return fmt.Sprintf("{block_height: %d, block_hash: %s}", block.Height, blockHash.String())
}

//GetHeadersMessage is one of the bytom msg type
type GetHeadersMessage struct {
	RawBlockLocator [][32]byte
	RawStopHash     [32]byte
}

//NewGetHeadersMessage return a new GetHeadersMessage
func NewGetHeadersMessage(blockLocator []*types.Hash, stopHash *types.Hash) *GetHeadersMessage {
	msg := &GetHeadersMessage{
		RawStopHash: stopHash.Value(),
	}
	for _, hash := range blockLocator {
		msg.RawBlockLocator = append(msg.RawBlockLocator, hash.Value())
	}
	return msg
}

//GetBlockLocator return the locator of the msg
func (m *GetHeadersMessage) GetBlockLocator() []*types.Hash {
	blockLocator := []*types.Hash{}
	for _, rawHash := range m.RawBlockLocator {
		hash := types.Hash(rawHash)
		blockLocator = append(blockLocator, &hash)
	}
	return blockLocator
}

func (m *GetHeadersMessage) String() string {
	return fmt.Sprintf("{stop_hash: %s}", hex.EncodeToString(m.RawStopHash[:]))
}

//GetStopHash return the stop hash of the msg
func (m *GetHeadersMessage) GetStopHash() *types.Hash {
	hash := types.Hash(m.RawStopHash)
	return &hash
}

//HeadersMessage is one of the bytom msg type
type HeadersMessage struct {
	RawHeaders [][]byte
}

//NewHeadersMessage create a new HeadersMessage
func NewHeadersMessage(headers []*types.BlockHeader) (*HeadersMessage, error) {
	RawHeaders := [][]byte{}
	for _, header := range headers {
		data, err := json.Marshal(header)
		if err != nil {
			return nil, err
		}

		RawHeaders = append(RawHeaders, data)
	}
	return &HeadersMessage{RawHeaders: RawHeaders}, nil
}

//GetHeaders return the headers in the msg
func (m *HeadersMessage) GetHeaders() ([]*types.BlockHeader, error) {
	headers := []*types.BlockHeader{}
	for _, data := range m.RawHeaders {
		header := &types.BlockHeader{}
		if err := json.Unmarshal(data, header); err != nil {
			return nil, err
		}

		headers = append(headers, header)
	}
	return headers, nil
}

func (m *HeadersMessage) String() string {
	return fmt.Sprintf("{header_length: %d}", len(m.RawHeaders))
}

//GetBlocksMessage is one of the bytom msg type
type GetBlocksMessage struct {
	RawBlockLocator [][32]byte
	RawStopHash     [32]byte
}

//NewGetBlocksMessage create a new GetBlocksMessage
func NewGetBlocksMessage(blockLocator []*types.Hash, stopHash *types.Hash) *GetBlocksMessage {
	msg := &GetBlocksMessage{
		RawStopHash: stopHash.Value(),
	}
	for _, hash := range blockLocator {
		msg.RawBlockLocator = append(msg.RawBlockLocator, hash.Value())
	}
	return msg
}

//GetBlockLocator return the locator of the msg
func (m *GetBlocksMessage) GetBlockLocator() []*types.Hash {
	blockLocator := []*types.Hash{}
	for _, rawHash := range m.RawBlockLocator {
		hash := types.Hash(rawHash)
		blockLocator = append(blockLocator, &hash)
	}
	return blockLocator
}

//GetStopHash return the stop hash of the msg
func (m *GetBlocksMessage) GetStopHash() *types.Hash {
	hash := types.Hash(m.RawStopHash)
	return &hash
}

func (m *GetBlocksMessage) String() string {
	return fmt.Sprintf("{stop_hash: %s}", hex.EncodeToString(m.RawStopHash[:]))
}

//BlocksMessage is one of the bytom msg type
type BlocksMessage struct {
	RawBlocks [][]byte
}

//NewBlocksMessage create a new BlocksMessage
func NewBlocksMessage(blocks []*types.Block) (*BlocksMessage, error) {
	rawBlocks := [][]byte{}
	for _, block := range blocks {
		data, err := json.Marshal(block)
		if err != nil {
			return nil, err
		}

		rawBlocks = append(rawBlocks, data)
	}
	return &BlocksMessage{RawBlocks: rawBlocks}, nil
}

//GetBlocks returns the blocks in the msg
func (m *BlocksMessage) GetBlocks() ([]*types.Block, error) {
	blocks := []*types.Block{}
	for _, data := range m.RawBlocks {
		block := &types.Block{}
		if err := json.Unmarshal(data, block); err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (m *BlocksMessage) String() string {
	return fmt.Sprintf("{blocks_length: %d}", len(m.RawBlocks))
}

//StatusRequestMessage status request msg
type StatusRequestMessage struct{}

func (m *StatusRequestMessage) String() string {
	return "{}"
}

//StatusResponseMessage get status response msg
type StatusResponseMessage struct {
	Height      uint64
	RawHash     [32]byte
	GenesisHash [32]byte
}

//NewStatusResponseMessage construct get status response msg
func NewStatusResponseMessage(blockHeader *types.BlockHeader, hash *types.Hash) *StatusResponseMessage {
	return &StatusResponseMessage{
		Height:      blockHeader.Height,
		RawHash:     blockHeader.Hash(),
		GenesisHash: hash.Value(),
	}
}

//GetHash get hash from msg
func (m *StatusResponseMessage) GetHash() *types.Hash {
	hash := types.Hash(m.RawHash)
	return &hash
}

//GetGenesisHash get hash from msg
func (m *StatusResponseMessage) GetGenesisHash() *types.Hash {
	hash := types.Hash(m.GenesisHash)
	return &hash
}

func (m *StatusResponseMessage) String() string {
	return fmt.Sprintf("{height: %d, hash: %s}", m.Height, hex.EncodeToString(m.RawHash[:]))
}

//TransactionMessage notify new tx msg
type TransactionMessage struct {
	RawTx []byte
}

//NewTransactionMessage construct notify new tx msg
func NewTransactionMessage(tx *types.Tx) (*TransactionMessage, error) {
	rawTx, err := tx.MarshalText()
	if err != nil {
		return nil, err
	}
	return &TransactionMessage{RawTx: rawTx}, nil
}

//GetTransaction get tx from msg
func (m *TransactionMessage) GetTransaction() (*types.Tx, error) {
	tx := &types.Tx{}
	if err := tx.UnmarshalText(m.RawTx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (m *TransactionMessage) String() string {
	tx, err := m.GetTransaction()
	if err != nil {
		return "{err: wrong message}"
	}
	return fmt.Sprintf("{tx_size: %d, tx_hash: %s}", len(m.RawTx), tx.Hash().String())
}

//MineBlockMessage new mined block msg
type MineBlockMessage struct {
	RawBlock []byte
}

//NewMinedBlockMessage construct new mined block msg
func NewMinedBlockMessage(block *types.Block) (*MineBlockMessage, error) {
	rawBlock, err := block.MarshalText()
	if err != nil {
		return nil, err
	}
	return &MineBlockMessage{RawBlock: rawBlock}, nil
}

//GetMineBlock get mine block from msg
func (m *MineBlockMessage) GetMineBlock() (*types.Block, error) {
	block := &types.Block{}
	if err := block.UnmarshalText(m.RawBlock); err != nil {
		return nil, err
	}
	return block, nil
}

func (m *MineBlockMessage) String() string {
	block, err := m.GetMineBlock()
	if err != nil {
		return "{err: wrong message}"
	}
	blockHash := block.Hash()
	return fmt.Sprintf("{block_height: %d, block_hash: %s}", block.Height, blockHash.String())
}

//FilterLoadMessage tells the receiving peer to filter the transactions according to address.
type FilterLoadMessage struct {
	Addresses [][]byte
}

func (m *FilterLoadMessage) String() string {
	return fmt.Sprintf("{addresses_length: %d}", len(m.Addresses))
}

// FilterAddMessage tells the receiving peer to add address to the filter.
type FilterAddMessage struct {
	Address []byte
}

func (m *FilterAddMessage) String() string {
	return fmt.Sprintf("{address: %s}", hex.EncodeToString(m.Address))
}

//FilterClearMessage tells the receiving peer to remove a previously-set filter.
type FilterClearMessage struct{}

func (m *FilterClearMessage) String() string {
	return "{}"
}
