package data_source

import (
	"fmt"
	"reflect"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slices"
)

const (
	GET_DATA_DEFAULT_LIMIT = 500
)

var SquirrelPgsql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Config struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DBName   string `json:"dbName" yaml:"dbName"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	SSLMode  string `json:"sslMode" yaml:"sslMode"`

	MaxOpenConnections    int `json:"maxOpenConnections" yaml:"maxOpenConnections"`
	MaxLifeTimeConnection int `json:"maxLifeTimeConnection" yaml:"maxLifeTimeConnection"` // Seconds
	MaxIdleConnections    int `json:"maxIdleConnections" yaml:"maxIdleConnections"`
	MaxIdleTimeConnection int `json:"maxIdleTimeConnection" yaml:"maxIdleTimeConnection"` // Seconds
}

// NewDB create new DB connection.
func NewDB(config *Config) (*sqlx.DB, error) {
	dsn := GetDsn(config)

	conn, err := sqlx.Open(config.Driver, dsn)
	if err != nil {
		return nil, err
	}

	conn.SetConnMaxLifetime(time.Duration(config.MaxLifeTimeConnection) * time.Second)
	conn.SetMaxOpenConns(config.MaxOpenConnections)
	conn.SetMaxIdleConns(config.MaxIdleConnections)
	conn.SetConnMaxIdleTime(time.Duration(config.MaxIdleTimeConnection) * time.Second)

	if er := conn.Ping(); er != nil {
		return nil, er
	}

	return conn, nil
}

func GetDsn(config *Config) string {
	if config.Driver == "mysql" {
		return fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?parseTime=true",
			config.User, config.Password, config.Host, config.Port, config.DBName,
		)
	} else if config.Driver == "postgres" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
		)
	}
	return ""
}

func GetDbColumnsAndValue(data any, excluded ...string) map[string]interface{} {
	v := reflect.ValueOf(data)
	typeOfS := v.Type()

	res := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		variableValue := v.Field(i).Interface()
		columName := typeOfS.Field(i).Tag.Get("db")
		if !slices.Contains(excluded, columName) {
			res[columName] = variableValue
		}
	}
	return res
}

func GetDbColumns(data any, excluded ...string) []string {
	v := reflect.ValueOf(data)
	typeOfS := v.Type()

	res := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		columName := typeOfS.Field(i).Tag.Get("db")
		if !slices.Contains(excluded, columName) {
			res = append(res, columName)
		}
	}
	return res
}
