module github.com/opentracing-contrib/gorm

require (
	cloud.google.com/go v0.40.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.10
	github.com/lib/pq v1.1.1
	github.com/opentracing-contrib/sql v0.0.0-20191016160643-5c2172400eb5
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	golang.org/x/crypto v0.1.0 // indirect
)

replace github.com/opentracing-contrib/sql => github.com/opentracing-lib/sql v0.0.0-20191016160643-5c2172400eb5

go 1.13
