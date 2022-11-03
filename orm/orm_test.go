package orm

import (
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gowins/dionysus/log"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

// TestDsnInfo test convert DsnInfo to string
func TestDsnInfo(t *testing.T) {
	convey.Convey("convert DsnInfo to string", t, func() {
		d := DsnInfo{
			User:   "user",
			Passwd: "pass",
			Host:   "127.0.0.1",
			Port:   3306,
			DbName: "dbname",
		}
		m := map[string][]string{
			"charset":   {"utf8mb4"},
			"parseTime": {"True"},
			"loc":       {"Local"},
		}
		d.Params = m
		convey.So(d.String(), convey.ShouldEqual, "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&loc=Local&parseTime=True")
		dialector := d.Dialector()
		convey.So(dialector.Name(), convey.ShouldEqual, "mysql")
	})
}

// TestSetUp orm set up
func TestSetUp(t *testing.T) {
	convey.Convey("db setup", t, func() {
		defer func() {
			_ = recover()
		}()
		d := DsnInfo{
			User:   "user",
			Passwd: "pass",
			Host:   "127.0.0.1",
			Port:   3306,
			DbName: "dbname",
		}
		log.Setup(log.SetProjectName("dbsetup"), log.WithWriter(io.Discard), log.WithOnFatal(&log.MockCheckWriteHook{}))
		Setup(d.Dialector())
	})
}

func TestSqlMock(t *testing.T) {
	convey.Convey("sql mock", t, func() {
		defer func() {
			_ = recover()
		}()
		db, mock, err := sqlmock.New()
		convey.So(err, convey.ShouldBeNil)
		defer db.Close()
		rows := sqlmock.NewRows([]string{"SELECT VERSION()"}).AddRow(1)
		mock.ExpectQuery("SELECT VERSION()").WillReturnRows(rows)
		dialector := DialectorByDB(db)
		log.Setup(log.SetProjectName("sqlmock"), log.WithWriter(io.Discard), log.WithOnFatal(&log.MockCheckWriteHook{}))
		Setup(dialector, testFnOpts(), testGormOpts())
		convey.So(GetDefaultDB(), convey.ShouldNotBeNil)
	})
}

func testGormOpts() ConfigOpt {
	return WithGormOpts(&gorm.Config{})
}

func testFnOpts() ConfigOpt {
	return WithOptFns(WithMaxOpenConns(10),
		WithMaxIdleConns(10),
		WithConnMaxIdleTime(time.Second*5),
		WithConnMaxLifetime(time.Second*5),
		WithCharset("utf8mb4"),
		WithParseTime("True"),
		WithLoc("local"),
	)
}
