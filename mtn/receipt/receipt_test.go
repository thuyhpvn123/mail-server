package receipt

import (
	"reflect"
	"testing"

	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
)

func TestReceipt_Json(t *testing.T) {
	type fields struct {
		proto *pb.Receipt
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestJson",
			fields: fields{
				proto: &pb.Receipt{
					TransactionHash: []byte{1, 2, 3},
					FromAddress:     []byte{1, 2, 3},
					ToAddress:       []byte{1, 2, 3},
					Amount:          []byte{1, 2, 3},
					Action:          pb.ACTION_EMPTY,
					Status:          pb.RECEIPT_STATUS_RETURNED,
					Return:          []byte{1, 2, 3},
					Exception:       pb.EXCEPTION_NONE,
					GasUsed:         10,
					GasFee:          10,
				},
			},
			want:    []byte(`{"action":0,"amount":"0x10203","exception":-1,"from_address":"0x0000000000000000000000000000000000010203","gas_fee":10,"gas_used":10,"return_value":"010203","status":0,"to_address":"0x0000000000000000000000000000000000010203","transaction_hash":"0x0000000000000000000000000000000000000000000000000000000000010203"}                              `),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Receipt{
				proto: tt.fields.proto,
			}
			got, err := r.Json()
			logger.Info("Got json", string(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("Receipt.Json() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Receipt.Json() = %v, want %v", got, tt.want)
			}
		})
	}
}
