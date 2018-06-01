package comm

import (
	"net/http"
	"reflect"
	"testing"
)

func TestPushData(t *testing.T) {
	//server side
	sk := "AAAAn1EFrO0:APA91bHB0iXK2rF-6i9UfGwDgKiUJMoN82zniuSSp1evXvoc34pMIpt08z2JHlxk9Z1h2xOcfrS8b5M8dFTT1eTWJpDUpKZ_oSHLFZyFxVqa5PmX11pgHf7ByRLwZs1pMfc_8VrjeB9w"
	cl := http.DefaultClient
	//client side
	ft := "cTHRtprTvK8:APA91bEfBpPVuobUFeUMFR_8-wpHn6nIFEqS5pOVXRY3xJOE9n5fxdDKGqzj-g_NHntWCywSpsoUx4Plk3Gr795rFldgUeRug2Y6mtT351IwN46LMUJWs_IhGLZ3SQOytYTsn5VEo77H"
	pr := []Pair{Pair{Key: "state", Value: "news"},
		Pair{Key: "type", Value: "openPage"},
		Pair{Key: "title", Value: "Título da Mensagem de Teste"},
		Pair{Key: "message", Value: "Mensagem de Teste com vários caracteres"}}
	//result
	rp := FCMResp{Code: 200, Success: true, ErrMsg: ""}
	//FCM_PLUGIN_ACTIVITY
	type fields struct {
		serverKey string
		client    *http.Client
	}
	type args struct {
		ftoken string
		pairs  []Pair
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   FCMResp
	}{
		{"teste_success_0", fields{sk, cl}, args{ft, pr}, rp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fcmNotifierInst{
				serverKey: tt.fields.serverKey,
				client:    tt.fields.client,
			}
			if got := f.PushData(tt.args.ftoken, "Title Test", "Msg Test", "icon", tt.args.pairs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fcmNotifierInst.PushData() = %v, want %v", got, tt.want)
			}
		})
	}
}
