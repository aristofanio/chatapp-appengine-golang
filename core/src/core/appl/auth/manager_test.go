package auth

import (
	"context"
	"core/infra/data"
	"core/infra/data/uuid"
	"reflect"
	"testing"

	"google.golang.org/appengine/aetest"
)

func Test_StoreAndCheckGuest(t *testing.T) {
	//
	ctx, done, _ := aetest.NewContext()
	defer done()
	//
	guid := uuid.NewUID("guest")
	gexp := &data.Guest{
		ID: guid,
	}
	//
	//
	type fields struct {
		ctx context.Context
	}
	type args struct {
		uGuestId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data.Guest
		wantErr bool
	}{
		{"test_success_0", fields{ctx}, args{"uGuestID_0"}, gexp, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//register an guest
			//TODO: register an Guest
			//create manager
			m := managerInst{
				ctx: tt.fields.ctx,
			}
			got, err := m.checkGuest(tt.args.uGuestId)
			if (err != nil) != tt.wantErr {
				t.Errorf("managerInst.checkGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("managerInst.checkGuest() = %v, want %v", got, tt.want)
			}
		})
	}
}
