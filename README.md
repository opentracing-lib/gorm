# gorm  <github.com/opentracing-contrib/gorm>


Package gorm provides a means of instrumenting GORM database operations.

To trace GORM operations, import the appropriate apmgorm/dialects package (instead of the gorm/dialects package), and use apmgorm.Open (instead of gorm.Open). The parameters are exactly the same.

Once you have a *gorm.DB from apmgorm.Open, you can call apmgorm.WithContext to propagate a context containing a transaction to the operations:

```
import (
    apmgorm "github.com/opentracing-contrib/gorm"
    _ "github.com/opentracing-contrib/gorm/dialects/mysql"
)

func main() {
	db, err := apmgorm.Open("mysql", "")
	...
	db = apmgorm.WithContext(ctx, db)
	db.Find(...) // creates a "SELECT FROM <foo>" span
}
```
