package main

import (
	"fmt"
	"io/ioutil"

	"github.com/msharran/labs/go-protobuff/config"
	"google.golang.org/protobuf/proto"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Print("generating config...")
	p := config.Person{
		Name:     "Sharran",
		EmailIds: []string{"sharran.murali@gmail.com"},
		Gender:   config.Gender_MALE,
		Tags: map[string]string{
			"foo": "bar",
		},
	}

	out, err := proto.Marshal(&p)
	check(err)

	ioutil.WriteFile("out", out, 0644)
	fmt.Println("done")

	fmt.Println("reading from proto file")

	in, err := ioutil.ReadFile("out")
	check(err)

	person := &config.Person{}
	check(proto.Unmarshal(in, person))

	fmt.Printf("%+v\n", person)
}
