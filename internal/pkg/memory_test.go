package memory

import (
	"io"
	"reflect"
	"testing"
)

func TestParseMemInfo(t *testing.T) {
	type args struct {
		memInfoSource io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    MemInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemInfo(tt.args.memInfoSource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMemInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMemInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
