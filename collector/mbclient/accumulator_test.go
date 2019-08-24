package mbclient

import (
	"reflect"
	"testing"

	"github.com/mythay/spider"
)

func Test_accumulate_no_range(t *testing.T) {
	device := spider.CfgDevice{
		Register: map[string]spider.CfgRegister{"value-1": {Base: 0, Type: "float"}, "value-2": {Base: 4, Type: "signed"}},
		Range:    nil,
	}
	type args struct {
		collect []string
		device  spider.CfgDevice
	}
	tests := []struct {
		name    string
		args    args
		want    []spider.CfgRange
		wantErr bool
	}{
		// TODO: Add test cases.
		{"only one float reg", args{[]string{"value-1"}, device}, []spider.CfgRange{{0, 2}}, false},
		{"only one singed reg", args{[]string{"value-2"}, device}, []spider.CfgRange{{4, 1}}, false},
		{"two regs", args{[]string{"value-1", "value-2"}, device}, []spider.CfgRange{{0, 2}, {4, 1}}, false},
		{"invalid value", args{[]string{"value-notexist"}, device}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aggregateRange(tt.args.collect, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("accumulate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accumulate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accumulate_with_range(t *testing.T) {
	device := spider.CfgDevice{
		Register: map[string]spider.CfgRegister{"value-1": {Base: 0, Type: "float"}, "value-2": {Base: 4, Type: "signed"}, "value-3": {Base: 6, Type: "signed"}},
		Range:    []spider.CfgRange{{0, 10}},
	}
	type args struct {
		collect []string
		device  spider.CfgDevice
	}
	tests := []struct {
		name    string
		args    args
		want    []spider.CfgRange
		wantErr bool
	}{
		// TODO: Add test cases.
		{"only one float reg", args{[]string{"value-1"}, device}, []spider.CfgRange{{0, 2}}, false},
		{"only one singed reg", args{[]string{"value-2"}, device}, []spider.CfgRange{{4, 1}}, false},
		{"two regs", args{[]string{"value-1", "value-2"}, device}, []spider.CfgRange{{0, 5}}, false},
		{"two regs, not start with 0", args{[]string{"value-2", "value-3"}, device}, []spider.CfgRange{{4, 3}}, false},
		{"invalid value", args{[]string{"value-notexist"}, device}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aggregateRange(tt.args.collect, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("accumulate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accumulate() = %v, want %v", got, tt.want)
			}
		})
	}
}
