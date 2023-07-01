package main

import (
	"fmt"
	"log"
	"os/exec"
)

func identifyJavaPid() string {
	app := "pgrep"
	arg0 := "-o"
	arg1 := "java"

	cmd := exec.Command(app, arg0, arg1)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("No Java PID found!")
	}
	return string(stdout)
}

func main() {

	fmt.Printf("Java PID identified: %s", identifyJavaPid())
}
