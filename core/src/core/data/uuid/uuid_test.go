package uuid

import (
	"fmt"
	"regexp"
	"testing"
)

func TestNewUID(t *testing.T) {
	//
	prfx := "test"
	//
	xprxSuccess := fmt.Sprintf("%s-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", prfx)
	wntSuccess, _ := regexp.Compile(xprxSuccess)
	//
	xprxFail := fmt.Sprintf("%s-[0-9a-f]{7}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", prfx)
	wntFail, _ := regexp.Compile(xprxFail)
	//
	tests := []struct {
		name string
		rgxp *regexp.Regexp
		want bool
	}{
		{"test_with_success_0", wntSuccess, true},
		{"test_with_fail_0", wntFail, false},
	}
	for _, tt := range tests {
		//
		t.Run(tt.name, func(t *testing.T) {
			//
			if got := NewUID(prfx); !(tt.rgxp.MatchString(string(got)) == tt.want) {
				t.Errorf("NewUID() = %v, want matching with %v", got, tt.rgxp)
			}
		})
	}
}
