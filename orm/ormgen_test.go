package orm

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gen"
)

func TestGormGen(t *testing.T) {
	convey.Convey("generate code", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		convey.So(err, convey.ShouldBeNil)
		defer db.Close()
		row := sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.6.16-log")
		mock.ExpectQuery("SELECT VERSION()").WillReturnRows(row)

		rowDatabase := sqlmock.NewRows([]string{"DATABASE()"}).AddRow("rds5ewin")
		mock.ExpectQuery("SELECT DATABASE()").WillReturnRows(rowDatabase)

		rowSchema := sqlmock.NewRows([]string{"SCHEMA_NAME "}).AddRow("rds5ewin")
		mock.ExpectQuery("SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE ? ORDER BY SCHEMA_NAME=? DESC,SCHEMA_NAME limit 1").
			WithArgs("rds5ewin%", "rds5ewin").
			WillReturnRows(rowSchema)

		rowsSelec := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "tomoki", "example@gmail.com", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT * FROM `users` LIMIT 1").WillReturnRows(rowsSelec)

		mock.ExpectQuery("SELECT column_name, column_default, is_nullable = 'YES', data_type, character_maximum_length, column_type, column_key, extra, column_comment, numeric_precision, numeric_scale , datetime_precision FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION").
			WithArgs("rds5ewin", "users").
			WillReturnRows(sqlRows())

		defaultOutPath = filepath.Join(os.TempDir(), "repository")
		defaultPkgPath = filepath.Join(os.TempDir(), "model")
		ormGen := OrmGenCmd{
			Dialector: DialectorByDB(db),
			Cfg:       gen.Config{},
			DataTyMap: map[string]func(string) string{
				"tinyint": func(s string) string { return "int8" },
			},
		}
		ormGen.GetShutdownFunc()()
		err = ormGen.GetCmd().Execute()
		convey.So(err, convey.ShouldNotBeNil)
		ormGen.GenModels = []GenModel{
			{TableName: "users"},
		}
		err = ormGen.GetCmd().Execute()
		convey.So(err, convey.ShouldBeNil)
		convey.So(mock.ExpectationsWereMet(), convey.ShouldBeNil)
		os.RemoveAll(defaultOutPath)
		os.RemoveAll(defaultPkgPath)
	})
}

func sqlRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"column_name", "column_default", "is_nullable = 'YES'", "data_type", "character_maximum_length", "column_type", "column_key", "extra", "column_comment", "numeric_precision", "numeric_scale", "datetime_precision"}).
		AddRow("id", 0, 0, "int", nil, "int(10) unsigned", "PRI", "", "primarykey", 10, 0, nil).
		AddRow("name", "", 0, "varchar", nil, "varchar(10)", "", "", "name", nil, 0, nil).
		AddRow("email", "", 0, "varchar", nil, "varchar(10)", "", "", "email", nil, 0, nil).
		AddRow("created_at", nil, 0, "datetime", nil, "datetime", "", "", "created_at", nil, 0, nil).
		AddRow("updated_at", nil, 0, "datetime", nil, "datetime", "", "", "updated_at", nil, 0, nil).
		AddRow("deleted_at", nil, 0, "datetime", nil, "datetime", "", "", "deleted_at", nil, 0, nil)
}
