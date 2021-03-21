package main

import (
	"fmt"
	"log"
	"time"
)

func Log(msg string, args ...string) {
	log.Printf(msg + "\n".args...)
}

const ticksPerSec int64 = 2

func getTick() float64 {
	now := time.Now().UnixNano()
	return float64(now) / 1000 / 1000

}

func main() {
	//skipTicks := 1000 / ticksPerSec
	//lastTick := getTick()
	//for {
	//lastTick += skipTicks
	//delta := lastTick - getTick()
	//if delta > 0 {
	//time.Sleep(time.Duration(delta) * time.Millisecond)
	//}

	////log.Println("update")
	//fmt.Println(skipTicks, lastTick, delta, delta > skipTicks)
	//time.Sleep(150 * time.Millisecond)
	//}

	//ticker := time.NewTicker(time.Duration(1000/ticksPerSec) * time.Millisecond)
	ticker := time.NewTicker(time.Second / time.Duration(ticksPerSec))
	defer ticker.Stop()
	last := time.Now()
	for {
		select {
		case t := <-ticker.C:
			fmt.Println(t, t.Sub(last))
			last = t
		}
		//time.Sleep(950 * time.Millisecond)
	}
}
