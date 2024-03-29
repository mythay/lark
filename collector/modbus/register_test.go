package modbus

import (
	"reflect"
	"testing"

	cfg "github.com/mythay/lark/config"
	"github.com/stretchr/testify/assert"
)

// func TestConnectModbusTcp(t *testing.T) {
// 	assert := assert.New(t)

// 	regs := []regReq{{3, 0, 2}}
// 	x := NewClient("127.0.0.1:502", regs)
// 	assert.NotNil(x)
// 	assert.Nil(x.Once())
// }

func Benchmark_Parse(b *testing.B) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		reg     *mbReg
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for i := 0; i < b.N; i++ {

		for _, tt := range tests {
			_, _ = tt.reg.parse(tt.args.data)
		}
	}
}

func Benchmark_ReflectParse(b *testing.B) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		reg     *mbReg
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for i := 0; i < b.N; i++ {

		for _, tt := range tests {
			_, _ = tt.reg.reflectParse(tt.args.data)
		}
	}
}

func TestCfgRegister_Parse(t *testing.T) {
	// assert := assert.New(t)
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		reg     *mbReg
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reg.parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CfgRegister.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// assert.EqualValues(tt.want, got)
				t.Errorf("CfgRegister.Parse() = %v(%s), want %v(%s)", got, reflect.TypeOf(got), tt.want, reflect.TypeOf(tt.want))
			}
		})
	}
}

func TestCfgRegister_ReflectParse(t *testing.T) {
	// assert := assert.New(t)
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		reg     *mbReg
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &mbReg{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reg.reflectParse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CfgRegister.NoReflectParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// assert.EqualValues(tt.want, got)
				t.Errorf("CfgRegister.NoReflectParse() = %v(%s), want %v(%s)", got, reflect.TypeOf(got), tt.want, reflect.TypeOf(tt.want))
			}
		})
	}
}

func Test_parseTagString(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
		{"one key", args{"a=b"}, map[string]string{"a": "b"}},
		{"two key", args{"a=b b=c"}, map[string]string{"a": "b", "b": "c"}},
		{"one key but two value", args{"a=b a=c"}, map[string]string{"a": "c"}},
		{"user & as split letter", args{"a=b&b=c"}, map[string]string{"a": "b", "b": "c"}},
		{"key without value", args{"ab="}, map[string]string{}},
		{"no equal", args{"ab"}, map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTagString(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_map_user(t *testing.T) {
	assert := assert.New(t)

	var a map[string]string
	assert.Nil(a)
	b := map[string]string{}
	assert.NotNil(b)
}
