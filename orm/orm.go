package orm

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gowins/dionysus/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbs                 = make(map[string]*gorm.DB)
	defaultMaxLifetime  = 7200 * time.Second // 单位  time.Second
	defaultMaxOpenConns = 1000               // 设置数据库连接池最大连接数
	defaultMaxIdleConns = 100                // 连接池最大允许的空闲连接数，如果空闲的连接数大于100，超过的连接会被连接池关闭
	defaultCharset      = "utf8mb4"
	defaultParseTime    = "True"
	defaultLocal        = "Local"
)

var rw sync.RWMutex

// Options connect database options
type Options struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	charset         string
	parseTime       string
	loc             string
}

// DsnInfo user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
type DsnInfo struct {
	User   string
	Passwd string
	Host   string
	Port   uint32
	DbName string
	Params url.Values
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
		c.optFns = append(c.optFns, optFn...)
	}
}

// WithGormOpts set Config gormOpts
func WithGormOpts(opts ...gorm.Option) ConfigOpt {
	return func(c *Config) {
		c.gormOpts = append(c.gormOpts, opts...)
	}
}

// DbMap multi-database maps
type DbMap struct {
	Dialector gorm.Dialector
	Opts      []ConfigOpt
}

func Setup(dbMaps map[string]DbMap) {
	for name, dbMap := range dbMaps {
		ormDB := getDB(dbMap.Dialector, dbMap.Opts...)
		if name == "" || ormDB == nil {
			log.Fatalf("init gorm.DB failure: %v", name)
		}
		rw.Lock()
		defer rw.Unlock()
		if _, ok := dbs[name]; ok {
			log.Fatalf("repeat database: %v", name)
		}
		dbs[name] = ormDB
	}
}

func getDB(dialector gorm.Dialector, opts ...ConfigOpt) *gorm.DB {
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

	gormDB, err := gorm.Open(dialector, c.gormOpts...)
	if err != nil {
		log.Fatalf("open gorm failed %s", err.Error())
	}

	db, err := gormDB.DB()
	if err != nil {
		log.Fatalf("get gorm db failed %s", err.Error())
	}

	db.SetConnMaxLifetime(options.connMaxLifetime)
	db.SetMaxIdleConns(options.maxIdleConns)
	db.SetMaxOpenConns(options.maxOpenConns)
	return gormDB
}

func GetDB(name string) *gorm.DB {
	rw.RLock()
	defer rw.RUnlock()
	return dbs[name]
}
