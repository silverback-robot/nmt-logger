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
		log.Fatal("NativeMemoryTracking not enabled for PID: " + pid)
	}

	match, _ := regexp.MatchString("NativeMemoryTracking", string(stdout))
	return match
}

func getNativeMemoryData(pid string) string {
	command := "jcmd"
	arg0 := pid
	arg1 := "VM.native_memory"
	arg2 := "summary"
	cmd := exec.Command(command, arg0, arg1, arg2)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal("Error getting Native Memory Tracking data for PID: " + pid)
	}

	return string(stdout)
}

func parseNMTData(rawData string) map[string]string {
	lookupStrings := [11]string{"Total", "Java Heap", "Class", "Thread", "Code", "GC", "Compiler", "Internal", "Symbol", "Native Memory Tracking ", "Arena Chunk"}

	parsedData := make(map[string]string)
	for _, k := range lookupStrings {

		// RegExp to identify all punctuation (p{P}) and mathematical symbols (p{S})
		reg, _ := regexp.Compile(`[\p{P}\p{S}]+`)

		leftIdx := strings.Index(rawData, k)
		rightIdx := strings.Index(string(rawData[leftIdx:]), "\n")

		value := string(rawData[leftIdx : leftIdx+rightIdx])

		splitString := strings.Fields(reg.ReplaceAllString(value, " "))
		value = splitString[len(splitString)-1]
		key := (strings.ReplaceAll(strings.ToUpper(strings.Trim(k, " ")), " ", "_"))

		// fmt.Println(key, value)
		parsedData[key] = value
	}
	return parsedData
}

func main() {
	javaPid := identifyJavaPid()
	fmt.Printf("Java PID identified: %s\n", javaPid)
	match := checkNMTEnabled(javaPid)
	if match {
		nmtRawData := getNativeMemoryData(javaPid)
		parsedData := parseNMTData(nmtRawData)
		fmt.Println(parsedData)
	}
}
