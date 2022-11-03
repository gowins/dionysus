package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/orm"
	"gorm.io/gen"
)

func main() {
	dsn := orm.DsnInfo{
		User:   "",
		Passwd: "",
		Host:   "",
		Port:   3306,
		DbName: "",
		Params: map[string][]string{
			"charset":   {"utf8mb4"},
			"parseTime": {"True"},
			"loc":       {"Local"},
		},
	}
	tmp := filepath.Join(os.TempDir(), "gorm_gen_1234567890")
	fmt.Println(tmp)
	ormGen := orm.OrmGenCmd{
		Dialector: dsn.Dialector(),
		Cfg: gen.Config{
			OutPath:      filepath.Join(tmp, "repository"),
			ModelPkgPath: filepath.Join(tmp, "model"),
		},
		GenModels: []orm.GenModel{
			{TableName: "yw_user"},
		},
	}
	ormGen.RegShutdownFunc(cmd.StopStep{
		StopFn: func() {
			fmt.Println("shutdown func")
		},
		StepName: "remove_test_file",
	})
	dio := dionysus.NewDio()
	dio.DioStart("orm", &ormGen)
	os.RemoveAll(tmp)
}
