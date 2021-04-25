package main

import (
	"fmt"
	"log"
	"time"
)

func pinger(c chan string) {
	for i := 0; ; i++ {
		c <- "ping"
	}
}

func ponger(c chan string) {
	for {
		c <- "pong"
	}
}

func receiver(c chan string) {
	for {
		msg := <-c
		log.Println(msg)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	channel := make(chan string)
	go pinger(channel)
	go ponger(channel)
	go receiver(channel)
	var input string
	fmt.Scanln(&input)
}
