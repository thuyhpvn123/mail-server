package checkpoint

import (
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	validator_types "gomail/cmd/validator/types"
	"gomail/pkg/block"
	"gomail/pkg/poh"
	pb "gomail/pkg/proto"
	"gomail/pkg/state"
	"gomail/pkg/storage"
	trie_package "gomail/pkg/trie"
	"gomail/types"
)

type CheckPoint struct {
	lastFullBlock      types.FullBlock
	thisLeaderSchedule types.LeaderSchedule
	nextLeaderSchedule types.LeaderSchedule
}

func NewCheckPoint(
	lastFullBlock types.FullBlock,
	thisLeaderSchedule types.LeaderSchedule,
	nextLeaderSchedule types.LeaderSchedule,
) validator_types.Checkpoint {
	return &CheckPoint{
		lastFullBlock:      lastFullBlock,
		thisLeaderSchedule: thisLeaderSchedule,
		nextLeaderSchedule: nextLeaderSchedule,
	}
}

func (cp *CheckPoint) Proto() protoreflect.ProtoMessage {
	pbCheckpoint := &pb.Checkpoint{
		LastBlock:          cp.lastFullBlock.Block().Proto().(*pb.Block),
		ThisLeaderSchedule: cp.thisLeaderSchedule.Proto().(*pb.LeaderSchedule),
		NextLeaderSchedule: cp.nextLeaderSchedule.Proto().(*pb.LeaderSchedule),
	}
	return pbCheckpoint
}

func (cp *CheckPoint) Marshal() ([]byte, error) {
	return proto.Marshal(cp.Proto())
}

func (cp *CheckPoint) Unmarshal(
	b []byte,
) error {
	pbCheckpoint := &pb.Checkpoint{}
	err := proto.Unmarshal(b, pbCheckpoint)
	if err != nil {
		return err
	}
	lastBlock := block.NewBlock(pbCheckpoint.LastBlock)
	cp.lastFullBlock = &block.FullBlock{}
	cp.lastFullBlock.FromProto(&pb.FullBlock{
		Block: lastBlock.Proto().(*pb.Block),
	})
	cp.thisLeaderSchedule = &poh.LeaderSchedule{}
	cp.thisLeaderSchedule.FromProto(pbCheckpoint.ThisLeaderSchedule)
	cp.nextLeaderSchedule = &poh.LeaderSchedule{}
	cp.nextLeaderSchedule.FromProto(pbCheckpoint.NextLeaderSchedule)

	return nil
}

func (cp *CheckPoint) LastFullBlock() types.FullBlock {
	return cp.lastFullBlock
}

func (cp *CheckPoint) ThisLeaderSchedule() types.LeaderSchedule {
	return cp.thisLeaderSchedule
}

func (cp *CheckPoint) NextLeaderSchedule() types.LeaderSchedule {
	return cp.nextLeaderSchedule
}

func (cp *CheckPoint) Save(savePath string) error {
	b, err := cp.Marshal()
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, b, 0644)
	return err
}

func (cp *CheckPoint) Load(savePath string) error {
	b, err := os.ReadFile(savePath)
	if err != nil {
		return err
	}
	err = cp.Unmarshal(b)
	return err
}

func (cp *CheckPoint) AccountStatesManager(
	accountStatesDbPath string,
	dbType string,
) (types.AccountStatesManager, error) {
	db, err := storage.LoadDb(
		accountStatesDbPath,
		dbType,
	)
	if err != nil {
		return nil, err
	}
	rootHash := cp.lastFullBlock.Block().AccountStatesRoot()
	trie, err := trie_package.New(rootHash, db)
	if err != nil {
		return nil, err
	}
	accountStatesManager := state.NewAccountStatesManager(trie, db)
	return accountStatesManager, nil
}
