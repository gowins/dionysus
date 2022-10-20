package xviper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func initDir(t *testing.T) (afero.Fs, string, string) {
	dir := "/etc/xviper"
	memFs := afero.NewMemMapFs()
	err := memFs.Mkdir(dir, 0o777)
	convey.So(err, convey.ShouldBeNil)
	fp := absFilePath(t, filepath.Join(dir, "config.yaml"))
	file, err := memFs.Create(fp)
	convey.So(err, convey.ShouldBeNil)
	_, err = file.Write([]byte(`key: value`))
	convey.So(err, convey.ShouldBeNil)
	file.Close()
	return memFs, dir, fp
}

// TestViperx test Viperx
func TestViperx(t *testing.T) {
	convey.Convey("get Viper", t, func() {
		fs, _, fp := initDir(t)
		os.Setenv(xviperEnv, fp)
		_, err := SetUp(nil, WithFns(func(v *viper.Viper) { v.SetFs(fs) }))
		convey.So(err, convey.ShouldBeNil)
		convey.So(viperX, convey.ShouldNotBeNil)
	})
}

// TestSetUp test SetUp
func TestSetUp(t *testing.T) {
	convey.Convey("set up", t, func() {
		memFs, _, fp := initDir(t)
		os.Setenv(xviperEnv, fp)
		convey.Convey("no set env", func() {
			_, err := SetUp(nil, WithEnvName("DIO_VIPER"))
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("set env", func() {
			v, err := SetUp(nil, WithEnvName(xviperEnv), WithFns(func(v *viper.Viper) { v.SetFs(memFs) }))
			convey.So(err, convey.ShouldBeNil)
			convey.So(v.GetString("key"), convey.ShouldEqual, "value")
		})
		convey.Convey("unmarshal", func() {
			type kv struct {
				Key string `mapstructure:"key"`
			}
			var s kv
			v, err := SetUp(&s, WithEnvName(xviperEnv), WithFns(func(v *viper.Viper) { v.SetFs(memFs) }))
			convey.So(err, convey.ShouldBeNil)
			convey.So(v.GetString("key"), convey.ShouldEqual, "value")
			convey.So(s.Key, convey.ShouldEqual, "value")
		})
	})
}

// TestConfigName set config name config type and config paths
func TestConfigName(t *testing.T) {
	convey.Convey("config name", t, func() {
		memFs, dir, _ := initDir(t)
		convey.Convey("empty file", func() {
			xv := NewXviper()
			xv.SetFs(memFs)
			xv.ConfigName("", "yaml", dir)
			convey.So(xv.ReadInConfig(), convey.ShouldBeNil)
		})
		convey.Convey("error file type", func() {
			xv := NewXviper()
			xv.SetFs(memFs)
			xv.ConfigName("c", "y", dir)
			convey.So(xv.ReadInConfig(), convey.ShouldNotBeNil)
		})
		convey.Convey("file not exist", func() {
			xv := NewXviper()
			xv.SetFs(memFs)
			xv.ConfigName("c", "yaml", dir)
			convey.So(xv.ReadInConfig(), convey.ShouldNotBeNil)
		})
		convey.Convey("success", func() {
			xv := NewXviper()
			xv.SetFs(memFs)
			xv.ConfigName("config", "yaml", dir)
			convey.So(xv.ReadInConfig(), convey.ShouldBeNil)
			convey.So(xv.GetString("key"), convey.ShouldEqual, "value")
		})
	})
}

// TestSetEnv set env
func TestSetEnv(t *testing.T) {
	convey.Convey("set automatic env", t, func() {
		convey.Convey("nil replacer", func() {
			xv := NewXviper()
			convey.So(os.Setenv("XVIPER_TT", "tt"), convey.ShouldBeNil)
			xv.SetEnv("XVIPER", nil)
			convey.So(xv.GetString("TT"), convey.ShouldEqual, "tt")
		})
		convey.Convey("replacer", func() {
			xv := NewXviper()
			convey.So(os.Setenv("XVIPER.TT", "tt"), convey.ShouldBeNil)
			xv.SetEnv("XVIPER", strings.NewReplacer(".", "_"))
			convey.So(xv.GetString("TT"), convey.ShouldEqual, "tt")
		})
	})
}

func absFilePath(t *testing.T, path string) string {
	t.Helper()

	s, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}

	return s
}
