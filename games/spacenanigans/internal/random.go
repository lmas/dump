package internal

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	// Don't care too much about the quality of the ingame RNG right now
	rand.Seed(time.Now().UnixNano())
}

func randomLine(path string) string {
	// Based on https://en.wikipedia.org/wiki/Reservoir_sampling
	f, err := os.Open(path)
	if err != nil {
		Log("randomLine(): %s", err)
	}
	defer f.Close()

	var line string
	var index int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		index += 1
		if line == "" {
			line = strings.TrimSpace(scanner.Text())
			continue
		}
		if rand.Intn(index) <= 1 {
			line = strings.TrimSpace(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		Log("randomLine(): %s", err)
	}
	return line
}

func RandomFirstName(gender string) string {
	switch gender {
	case "male":
		return randomLine("data/names/first_male.txt")
	case "female":
		return randomLine("data/names/first_female.txt")
	default:
		return randomLine("data/names/first.txt")
	}
}

func RandomLastName() string {
	return randomLine("data/names/last.txt")
}
