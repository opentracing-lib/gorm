package mysql

import (
	"github.com/go-sql-driver/mysql"

	apmsql "github.com/opentracing-contrib/sql"
)

func init() {
	apmsql.Register("mysql", &mysql.MySQLDriver{}, apmsql.WithDSNParser(ParseDSN))
}
