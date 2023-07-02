package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func identifyJavaPid() string {
	command := "pgrep"
	arg0 := "-o"
	arg1 := "java"

	cmd := exec.Command(command, arg0, arg1)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("No Java PID found!")
	}
	return strings.ReplaceAll(string(stdout), "\n", "")
}

func checkNMTEnabled(pid string) bool {

	command := "ps"
	arg0 := "-p"
	arg1 := pid
	arg2 := "-o"
	arg3 := "args"
	arg4 := "--no-headers"

	cmd := exec.Command(command, arg0, arg1, arg2, arg3, arg4)

	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("NativeMemoryTracking not enabled for PID: " + pid)
	}

	match, _ := regexp.MatchString("NativeMemoryTracking", string(stdout))
	return match
}

func main() {
	javaPid := identifyJavaPid()
	fmt.Printf("Java PID identified: %s", javaPid)
	match := checkNMTEnabled(javaPid)
	fmt.Println(match)
}
