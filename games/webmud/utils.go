package main

import "log"

const gLogo = `
                 /##      /##   Welcome To...
                | ###    /###
                | ####  /####  /######   /######  /#######        /###     /######
                | ## ##/## ## /##__  ## /##__  ##| ##__  ##       |###    /########
                | ##  ###| ##| ##  \ ##| ##  \ ##| ##  \ ##     /#####   /###__  ###
                | ##\  # | ##| ##  | ##| ##  | ##| ##  | ##    /######  /###/  | ###
                | ## \/  | ##|  ######/|  ######/| ##  | ##   |__  ###  |__/   | ###
                |__/     |__/ \______/  \______/ |__/  |__/      | ###         \ ###
                                                                 | ###      /######/
                                                                 | ###     /#######
                                                                 | ###    |____  ###
                   /#######                                      | ###         | ###
                  | ##__  ##                                     | ###   /###  | ###
                  | ##  \ ##  /######   /#######  /######        | ###  | ###  \ ###
                  | #######  |____  ## /##_____/ /##__  ##      /#######| #########/
                  | ##__  ##  /#######|  ###### | ########     | #######|  #######/
                  | ##  \ ## /##__  ## \____  ##| ##_____/     |_______/ \______/
                  | #######/|  ####### /#######/|  #######
                  |_______/  \_______/|_______/  \_______/     <mark>Version 0.1</mark>
`

func Log(msg string, args ...interface{}) {
	log.Printf(msg+"\n", args...)
}

////////////////////////////////////////////////////////////////////////////////

// Mostly stolen from: https://github.com/gobwas/ws-examples/blob/master/src/gopool/pool.go

type WorkPool struct {
	queue chan func()
}

func NewWorkPool(size, workers int) *WorkPool {
	wp := &WorkPool{
		queue: make(chan func(), size),
	}

	for w := 0; w < workers; w++ {
		go func() {
			for work := range wp.queue {
				work()
			}
		}()
	}

	return wp
}

func (wp *WorkPool) Do(work func()) {
	// if the queue is full, we'll be blocked here until a slot is free
	wp.queue <- work
}

//func (wp *WorkPool) DoTimeout(work func(), timeout time.Duration) error {}
