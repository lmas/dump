package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// The time format to use
const time_format = "2006-01-02 15:04:05 MST"

// Run an infinite loop and print the current local time each minute.
func main() {
	for {
		t := time.Now()
		fmt.Println(t.Format(time_format), battery())
		time.Sleep(1 * time.Minute)
	}
}

func battery() string {
	out, err := exec.Command("acpi", "-b").Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
