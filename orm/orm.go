package orm

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gowins/dionysus/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	defaultDB           *gorm.DB
	defaultMaxLifetime  = 2 * time.Second // 单位  time.Second
	defaultMaxOpenConns = 100             // 设置数据库连接池最大连接数
	defaultMaxIdleConns = 5               // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于5，超过的连接会被连接池关闭
	defaultCharset      = "utf8mb4"
	defaultParseTime    = "True"
	defaultLocal        = "Local"
)

// Options connect database options
type Options struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	charset         string
	parseTime       string
	loc             string
}

// DriverTy driver type
type DriverTy string

const (
	Mysql     DriverTy = "mysql"
	Postgres  DriverTy = "postgres"
	Sqlite    DriverTy = "sqlite"
	SqlServer DriverTy = "sqlserver"
)

// DsnInfo user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
type DsnInfo struct {
	User   string
	Passwd string
	Host   string
	Port   uint32
	DbName string
	Params url.Values
	Driver DriverTy
}

// String convert DsnInfo to string
// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&loc=Local&parseTime=True
func (dsn *DsnInfo) String() string {
	dsnStr := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v", dsn.User, dsn.Passwd, dsn.Host, dsn.Port, dsn.DbName)
	if encode := dsn.Params.Encode(); encode != "" {
		dsnStr = strings.Join([]string{dsnStr, encode}, "?")
	}
	return dsnStr
}

// Dialector return gorm.Dialector
func (dsn *DsnInfo) Dialector() gorm.Dialector {
	return mysql.Open(dsn.String())
}

// DialectorByDB return gorm.Dialector by sql.DB
func DialectorByDB(sqlDB *sql.DB) gorm.Dialector {
	return mysql.New(mysql.Config{Conn: sqlDB})
}

// OptionFunc set Options function
type OptionFunc func(opts *Options)

// WithMaxOpenConns set maxOpenConns
func WithMaxOpenConns(maxOpenConns int) OptionFunc {
	return func(opts *Options) {
		opts.maxOpenConns = maxOpenConns
	}
}

// WithMaxIdleConns set maxIdleConns
func WithMaxIdleConns(maxIdleConns int) OptionFunc {
	return func(opts *Options) {
		opts.maxIdleConns = maxIdleConns
	}
}

// WithConnMaxLifetime set connMaxLifetime
func WithConnMaxLifetime(connmaxLifetime time.Duration) OptionFunc {
	return func(opts *Options) {
		opts.connMaxLifetime = connmaxLifetime
	}
}

// WithConnMaxIdleTime sett connMaxIdleTime
func WithConnMaxIdleTime(connMaxIdleTime time.Duration) OptionFunc {
	return func(opts *Options) {
		opts.connMaxIdleTime = connMaxIdleTime
	}
}

// WithCharset set charset
func WithCharset(charset string) OptionFunc {
	return func(opts *Options) {
		opts.charset = charset
	}
}

// WithParseTime set parseTime
func WithParseTime(parseTime string) OptionFunc {
	return func(opts *Options) {
		opts.parseTime = parseTime
	}
}

// WithLoc set loc
func WithLoc(loc string) OptionFunc {
	return func(opts *Options) {
		opts.loc = loc
	}
}

// Config contain database connect options and gorm options
type Config struct {
	optFns   []OptionFunc
	gormOpts []gorm.Option
}

// ConfigOpt set Config function
type ConfigOpt func(*Config)

// WithOptFns set Config optFns
func WithOptFns(optFn ...OptionFunc) ConfigOpt {
	return func(c *Config) {
		c.optFns = optFn
	}
}

// WithGormOpts set Config gormOpts
func WithGormOpts(opts ...gorm.Option) ConfigOpt {
	return func(c *Config) {
		c.gormOpts = opts
	}
}

func Setup(dialector gorm.Dialector, opts ...ConfigOpt) error {
	var err error
	options := &Options{
		maxOpenConns:    defaultMaxOpenConns,
		maxIdleConns:    defaultMaxIdleConns,
		connMaxLifetime: defaultMaxLifetime,
		charset:         defaultCharset,
		parseTime:       defaultParseTime,
		loc:             defaultLocal,
	}
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	for _, optFn := range c.optFns {
		optFn(options)
	}

	defaultDB, err = gorm.Open(dialector, c.gormOpts...)
	if err != nil {
		log.Fatalf("open gorm failed %s", err.Error())
		return err
	}

	db, err := defaultDB.DB()
	if err != nil {
		log.Fatalf("get gorm db failed %s", err.Error())
		return err
	}

	db.SetConnMaxLifetime(options.connMaxLifetime)
	db.SetMaxIdleConns(options.maxIdleConns)
	db.SetMaxOpenConns(options.maxOpenConns)
	return nil
}

func GetDefaultDB() *gorm.DB {
	return defaultDB
}
