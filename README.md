# gorm-logger-logrus

Logrus logger for gorm v2

```go
package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "github.com/nekomeowww/gorm-logger-logrus"
)

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
        Logger: gormloggerlogrus.New(gormloggerlogrus.Options{
            Logger:                    logrus.NewEntry(logrus.New()),
            LogLevel:                  gorm.Slient,
            IgnoreRecordNotFoundError: false,
            SlowThreshold:             time.Millisecond * 200,
            FileWithLineNumField:      "file",
        }),
    })
    if err != nil {
        panic("failed to connect database")
    }
}
```
