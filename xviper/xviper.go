package xviper

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var (
	configFile = "/etc/conf/config.yaml"
	xviperEnv  = "DIO_XVIPER_CONF"
)

// ViperX get *viper.Viper
func ViperX() *viper.Viper {
	return viper.GetViper()
}

// ViperFunc 设置viper.Viper属性函数
type ViperFn func(*viper.Viper)

// ViperOpt viper.Viper配置项
type ViperOpt struct {
	envName     string
	opts        []viper.Option
	fns         []ViperFn
	decoderOpts []viper.DecoderConfigOption
}

// ViperOptFn 设置ViperOpt属性
type ViperOptFn func(*ViperOpt)

// WithEnvName 设置环境变量
func WithEnvName(envName string) ViperOptFn {
	return func(vo *ViperOpt) {
		vo.envName = envName
	}
}

// WithOpts 设置viper.Option
func WithOpts(opts ...viper.Option) ViperOptFn {
	return func(vo *ViperOpt) {
		vo.opts = opts
	}
}

// WithFns 设置ViperFn
func WithFns(fns ...ViperFn) ViperOptFn {
	return func(vo *ViperOpt) {
		vo.fns = fns
	}
}

// WithDecoderOpts 设置viper.DecoderConfigOption
func WithDecoderOpts(opts ...viper.DecoderConfigOption) ViperOptFn {
	return func(vo *ViperOpt) {
		vo.decoderOpts = opts
	}
}

// SetUp 根据环境变量配置获取配置文件路径，viper.Viper读取配置文件中的内容加载到内存中
//
// DIO_XVIPER_CONF = /etc/config/config.toml
//
// 若rawVal不为空，则会把map[string]any解析到结构体中，例如
//
//	type Example struct {
//		Name string `mapstructure:"name"`
//		Age int `mapstructure:"int"`
//	}
//
// [dependencies]
//
// kio = "1.0"
//
// 解析toml文件，设置上述环境变量DIO_XVIPER_CONF，调用如下代码
//
//	func ExampleSetup_k() {
//		_, err := xviper.SetUp(nil, xviper.WithOpts(viper.KeyDelimiter("_")))
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(viper.GetString("dependencies_kio"))
//	}
func SetUp(rawVal any, opts ...ViperOptFn) (*viper.Viper, error) {
	vOpt := &ViperOpt{}
	for _, opt := range opts {
		opt(vOpt)
	}
	v := viper.NewWithOptions(vOpt.opts...)
	if len(vOpt.fns) > 0 {
		for _, fn := range vOpt.fns {
			fn(v)
		}
	}
	if vOpt.envName != "" {
		xviperEnv = vOpt.envName
	}
	if name := os.Getenv(xviperEnv); name != "" {
		configFile = name
	}
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	if rawVal != nil {
		if err := v.Unmarshal(rawVal, vOpt.decoderOpts...); err != nil {
			return nil, err
		}
	}
	ov := viper.GetViper()
	*ov = *v
	return v, nil
}

// Xviper viper.Viper wrapper
type Xviper struct {
	*viper.Viper
}

// NewXviper ...
func NewXviper(opts ...viper.Option) *Xviper {
	v := viper.NewWithOptions(opts...)
	return &Xviper{Viper: v}
}

// ConfigName according to file name, file type, and file paths set Viper
func (xv *Xviper) ConfigName(name, ty string, paths ...string) {
	xv.SetConfigName(name)
	xv.SetConfigType(ty)
	xv.AddConfigPath(".")
	for _, path := range paths {
		xv.AddConfigPath(path)
	}
}

// SetEnv set env
// r := strings.NewReplacer(".", "_", "-", "_")
// . -> _
// - -> _
func (xv *Xviper) SetEnv(envPrefix string, r *strings.Replacer, bindEnvs ...string) {
	xv.AutomaticEnv()
	xv.SetEnvPrefix(strings.ToUpper(envPrefix))
	if r != nil {
		xv.SetEnvKeyReplacer(r)
	}
	if len(bindEnvs) > 0 {
		xv.BindEnv(bindEnvs...)
	}
}
