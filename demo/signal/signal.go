package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	// Block until a signal is received.
	fmt.Println("Got signal:")
	go func() {
		s := <-c
		fmt.Println("Got signal:", s)
	}()
	time.Sleep(time.Minute)
}
