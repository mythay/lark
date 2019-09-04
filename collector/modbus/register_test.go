package modbus

import (
	"reflect"
	"testing"

	cfg "github.com/mythay/lark/config"
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
		reg     *regValue
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for i := 0; i < b.N; i++ {

		for _, tt := range tests {
			_, _ = tt.reg.Parse(tt.args.data)
		}
	}
}

func Benchmark_ReflectParse(b *testing.B) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		reg     *regValue
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for i := 0; i < b.N; i++ {

		for _, tt := range tests {
			_, _ = tt.reg.ReflectParse(tt.args.data)
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
		reg     *regValue
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reg.Parse(tt.args.data)
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
		reg     *regValue
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{00, 01}}, int16(1), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{00, 01}}, uint16(1), false},
		{"uint16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "uint16"}}, args{[]byte{0xff, 01}}, uint16(65281), false},

		{"int32", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32"}}, args{[]byte{00, 01, 00, 02}}, int32(131073), false},
		{"int32 inverse", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int32", Inverse: true}}, args{[]byte{00, 01, 00, 02}}, int32(65538), false},
		{"int16", &regValue{CfgRegister: &cfg.CfgRegister{Type: "int16"}}, args{[]byte{0xff, 01}}, int16(-255), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reg.ReflectParse(tt.args.data)
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
