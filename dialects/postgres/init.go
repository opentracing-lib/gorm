package pq

import (
	pq "github.com/lib/pq"

	apmsql "github.com/opentracing-contrib/sql"
)

func init() {
	apmsql.Register("postgres", &pq.Driver{}, apmsql.WithDSNParser(ParseDSN))
}
