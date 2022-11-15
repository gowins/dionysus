### orm使用

#### orm

```go
package main

import (
	"os"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/orm"
	"gorm.io/gen"
)

func main() {
	log.Setup(log.SetProjectName("dbsetup"), log.WithWriter(io.Stdout))
	d := orm.DsnInfo{
        // 用户名
		User:   "user",
        // 密码
		Passwd: "pass",
        // ip地址
		Host:   "127.0.0.1",
        // 端口
		Port:   3306,
        // 数据库名
		DbName: "dbname",
	}
	d1 := orm.DsnInfo{
        // 用户名
		User:   "user",
        // 密码
		Passwd: "pass",
        // ip地址
		Host:   "127.0.0.1",
        // 端口
		Port:   3306,
        // 数据库名
		DbName: "dbname",
	}
	m := map[string]orm.DbMap{
		"default": {Dialector: d.Dialector(), Opts: []orm.ConfigOpt{testFnOpts(), testGormOpts()}},
		"d1": {Dialector: d1.Dialector(), Opts: nil},
	}
	orm.Setup(m)
    orm.GetDB("default")
}

    
// testGormOpts 设置gorm.Option
func testGormOpts() ConfigOpt {
	return orm.WithGormOpts(&gorm.Config{})
}

// testFnOpts 设置sql.DB连接参数
func testFnOpts() ConfigOpt {
	return orm.WithOptFns(
        // 最大打开连接数
        orm.WithMaxOpenConns(10),
        // 最大空闲连接数
		orm.WithMaxIdleConns(10),
        // 连接有效期
		orm.WithConnMaxLifetime(time.Second*5),
        // 参开: https://github.com/go-sql-driver/mysql#dsn-data-source-name
        // 连接字符集
		orm.WithCharset("utf8mb4"),
        // 是否解析time.Time
		orm.WithParseTime("True"),
        // time.Time市区
		orm.WithLoc("local"),
	)
}
```

#### gorm_gen

```go
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
		User:   "user",
		Passwd: "passwd",
		Host:   "127.0.0.1",
		Port:   3306,
		DbName: "rds5ewin",
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
    // 移除生成代码
	os.RemoveAll(tmp)
}
```

运行

```
go run main.go gorm_gen
```
