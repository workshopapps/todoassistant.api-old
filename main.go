package main

import (
	_ "github.com/go-sql-driver/mysql"
	"test-va/cmd"
)

func main() {
	cmd.Setup()
}
