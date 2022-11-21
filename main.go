package main

import (
	"test-va/cmd"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cmd.Setup()
}


