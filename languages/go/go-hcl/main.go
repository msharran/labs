package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Config struct {
	Aws Aws `hcl:"aws,block"`
}

type Aws struct {
	Region string   `hcl:"region"`
	Other  hcl.Body `hcl:"other,remain"`
}

func main() {
	parser := hclparse.NewParser()
	f, diag := parser.ParseHCLFile("config.hcl")

	fmt.Println(diag)

	c := Config{}
	fmt.Println(gohcl.DecodeBody(f.Body, nil, &c))
	fmt.Println(c.Aws.Region)
}
