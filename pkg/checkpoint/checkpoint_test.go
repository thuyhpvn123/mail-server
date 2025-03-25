package checkpoint

import (
	"testing"

	"gomail/pkg/logger"
	"gomail/types"
)

func TestCheckPoint_Load(t *testing.T) {
	type fields struct {
		lastFullBlock      types.FullBlock
		thisLeaderSchedule types.LeaderSchedule
		nextLeaderSchedule types.LeaderSchedule
	}
	type args struct {
		savePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test load checkpoint",
			args: args{
				savePath: "check_point.dat",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &CheckPoint{
				lastFullBlock:      tt.fields.lastFullBlock,
				thisLeaderSchedule: tt.fields.thisLeaderSchedule,
				nextLeaderSchedule: tt.fields.nextLeaderSchedule,
			}
			if err := cp.Load(tt.args.savePath); (err != nil) != tt.wantErr {
				t.Errorf("CheckPoint.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			logger.Info(
				"Checkpoint account state root",
				cp.lastFullBlock.Block().AccountStatesRoot(),
			)
		})
	}
}
