package block

import (
	"fmt"
	"os"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	"gomail/mtn/receipt"
	"gomail/mtn/types"
)

type ConfirmedBlockData struct {
	header          types.BlockHeader
	receipts        []types.Receipt
	branchStateRoot e_common.Hash
	validatorSigns  map[e_common.Address][]byte
}

func NewConfirmedBlockData(
	header types.BlockHeader,
	receipts []types.Receipt,
	branchStateRoot e_common.Hash,
	validatorSigns map[e_common.Address][]byte,
) *ConfirmedBlockData {
	return &ConfirmedBlockData{
		header:          header,
		receipts:        receipts,
		branchStateRoot: branchStateRoot,
		validatorSigns:  validatorSigns,
	}
}

func (c *ConfirmedBlockData) Header() types.BlockHeader {
	return c.header
}

func (c *ConfirmedBlockData) Receipts() []types.Receipt {
	return c.receipts
}

func (c *ConfirmedBlockData) BranchStateRoot() e_common.Hash {
	return c.branchStateRoot
}

func (c *ConfirmedBlockData) ValidatorSigns() map[e_common.Address][]byte {
	return c.validatorSigns
}

func (c *ConfirmedBlockData) SetHeader(header types.BlockHeader) {
	c.header = header
}

func (c *ConfirmedBlockData) SetBranchStateRoot(rootHash e_common.Hash) {
	c.branchStateRoot = rootHash
}

func (c *ConfirmedBlockData) Proto() *pb.ConfirmedBlockData {
	validatorSigns := make(map[string][]byte)
	for k, v := range c.validatorSigns {
		validatorSigns[k.Hex()] = v
	}
	return &pb.ConfirmedBlockData{
		Header:          c.header.Proto(),
		Receipts:        receipt.ReceiptsToProto(c.receipts),
		BranchStateRoot: c.branchStateRoot.Bytes(),
		ValidatorSigns:  validatorSigns,
	}
}

func (c *ConfirmedBlockData) FromProto(pbData *pb.ConfirmedBlockData) {
	c.header = &BlockHeader{}
	c.header.FromProto(pbData.Header)
	c.receipts = receipt.ProtoToReceipts(pbData.Receipts)
	c.branchStateRoot = e_common.BytesToHash(pbData.BranchStateRoot)
	c.validatorSigns = make(map[e_common.Address][]byte)
	for k, v := range pbData.ValidatorSigns {
		c.validatorSigns[e_common.HexToAddress(k)] = v
	}
}

func (c *ConfirmedBlockData) Marshal() ([]byte, error) {
	return proto.Marshal(c.Proto())
}

func (c *ConfirmedBlockData) Unmarshal(cData []byte) error {
	pbData := &pb.ConfirmedBlockData{}
	err := proto.Unmarshal(cData, pbData)
	if err != nil {
		return err
	}
	c.FromProto(pbData)
	return nil
}

func LoadConfirmedBlockDataFromFile(path string) (types.ConfirmedBlockData, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cData := &ConfirmedBlockData{}
	err = cData.Unmarshal(raw)
	if err != nil {
		logger.DebugP("HERE")
		return nil, err
	}
	return cData, nil
}

func SaveConfirmedBlockDataToFile(cData types.ConfirmedBlockData, path string) error {
	// Ensure the directory exists
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logger.Error("Failed to create directory", err)
	}
	raw, err := cData.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(
		fmt.Sprintf("%v%v", path, "confirmed_block_data.dat"),
		raw,
		0644,
	)
}
