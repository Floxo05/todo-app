package main

import (
	"github.com/floxo05/todoapi/internal/tools"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	tools.DoMigration("down")
}
