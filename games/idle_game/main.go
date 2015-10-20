package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const BANNER = "Gone v0.1"

type Resource struct {
	Owns            int   // How much the player owns
	GatherTime      int64 // Time to gather
	GatherMaxAmount int   // Units recieved
}

var Resources = map[string]*Resource{
	"Wood":  &Resource{0, 1, 1},
	"Stone": &Resource{0, 1, 1},
}

type Building struct {
	Name   string
	Amount string
}

var Buildings = []Building{}

type Command struct {
	Name string
	Help string
	Func func(string)
}

var Commands = []Command{
	Command{"stats", "Show stats", cmd_stats},
	Command{"wood", "Gather wood", cmd_wood},
}

func main() {
	fmt.Println(BANNER)
	rand.Seed(1) // TODO: at least set it to the current time

	bfin := bufio.NewReader(os.Stdin)
	for {
		// Read input from the player
		fmt.Printf("\n>")
		input, err := bfin.ReadString('\n')
		if err == io.EOF {
			return
		}
		check_error(err)
		input = strings.TrimSpace(input)
		if len(input) < 1 {
			continue
		}
		parts := strings.SplitN(input, " ", 2)

		// Check if input matches a command and if so, run it
		for _, cmd := range Commands {
			if parts[0] == cmd.Name {
				if len(parts) == 2 {
					cmd.Func(parts[1])
				} else {
					cmd.Func("")
				}
				continue
			}
		}
	}
}

func check_error(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func cmd_stats(args string) {
	fmt.Println("Your resources:")
	for k, v := range Resources {
		fmt.Printf("%s\t%d (time %d, max %d)\n", k, v.Owns, v.GatherTime, v.GatherMaxAmount)
	}
}

func cmd_wood(args string) {
	wood := Resources["Wood"]
	fmt.Print("Gathering wood... ")
	time.Sleep(time.Duration(wood.GatherTime) * time.Second)

	w := rand.Intn(wood.GatherMaxAmount + 1)
	fmt.Printf("Found %d logs (%d).\n", w, wood.GatherMaxAmount)

	wood.Owns += w
}
