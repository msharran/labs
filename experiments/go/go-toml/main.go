package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Watcher Watcher `toml:"watcher"`
}

type Watcher struct {
	Queue   string                   `toml:"queue"`
	Aws     Aws                      `toml:"aws"`
	Tenants []map[string]interface{} `toml:"tenants"`
}

type Aws struct {
	AccountName string `toml:"account_name"`
	Region      string `toml:"region"`
}

func main() {
	var conf Config
	md, err := toml.DecodeFile("./config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Metadata")
	fmt.Println(md)
	fmt.Println("data")
	fmt.Println(conf)
}
