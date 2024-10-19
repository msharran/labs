package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

func main() {
	conf, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	fmt.Println("Port:", conf.Section("server").Key("port").MustInt())
	fmt.Println("Localhost:", conf.Section("server").Key("localhost").MustBool())
	fmt.Printf("logLevel = %+v\n", conf.Section("server").Key("logLevel").String())
	fmt.Printf("sqlitePath = %+v\n", conf.Section("sqlitedb").Key("path").String())
}
