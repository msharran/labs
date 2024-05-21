package main

import "os/exec"

func main() {
	// run ls using exec and print out to stdout
	cmd := exec.Command("ls")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	print(string(out))
}
