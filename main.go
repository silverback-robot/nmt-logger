package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func identifyJavaPid() string {
	command := "pgrep"
	arg0 := "-o"
	arg1 := "java"

	cmd := exec.Command(command, arg0, arg1)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("No Java PID found!")
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
		fmt.Println("NativeMemoryTracking not enabled for PID: " + pid)
		return false
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

func writeToFile(parsedData map[string]string) {
	// Check if log file exists; create file with header if not exists
	if _, err := os.Stat("logs/memory_stats.log"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("logs/", 0750)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		err = os.WriteFile("logs/memory_stats.log", []byte("ARENA_CHUNK|CLASS|CODE|COMPILER|GC|INTERNAL|JAVA_HEAP|NATIVE_MEMORY_TRACKING|SYMBOL|THREAD|TOTAL\n"), 0660)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Append parsed nmt values to log file
	f, err := os.OpenFile("logs/memory_stats.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(parsedData["ARENA_CHUNK"] + "|" + parsedData["CLASS"] + "|" + parsedData["CODE"] + "|" + parsedData["COMPILER"] + "|" + parsedData["GC"] + "|" + parsedData["INTERNAL"] + "|" + parsedData["JAVA_HEAP"] + "|" + parsedData["NATIVE_MEMORY_TRACKING"] + "|" + parsedData["SYMBOL"] + "|" + parsedData["THREAD"] + "|" + parsedData["TOTAL"] + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for {

		javaPid := identifyJavaPid()
		if javaPid != "" {
			fmt.Printf("Java PID identified: %s\n", javaPid)
		}
		match := checkNMTEnabled(javaPid)
		if match {
			nmtRawData := getNativeMemoryData(javaPid)
			parsedData := parseNMTData(nmtRawData)
			fmt.Println(parsedData)
			writeToFile(parsedData)
		}
		time.Sleep(60 * time.Second)
	}
}
