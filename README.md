<h2 align="center">GORM SQLCommenter Plugin</h2>

This plugin adds a new clause on GORM statement builder, allowing you to add tags to your requests using the [SQLCommenter](https://google.github.io/sqlcommenter/) protocol.

### How to use

```golang
package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  // Add plugin package
  sqlCommenter "github.com/gouyelliot/gorm-sqlcommenter-plugin"
)

func main() {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Register the plugin in GORM
    db.Use(sqlCommenter.New())

    // Add SQLCommenter tag to your request
    db.Clauses(sqlCommenter.NewTag("application", "api")).Table("test").Scan(nil)
}
```