package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type listRow struct {
	row string
}

type messageActor struct{}

var receiveCounter int
var sentCounter int

func (state *messageActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *listRow:
		receiveCounter++
		if receiveCounter%50000 == 0 {
			splitted := strings.Split(msg.row, ",")
			fmt.Println(receiveCounter, splitted[7])
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Create an actor
	props := actor.FromProducer(func() actor.Actor { return &messageActor{} })
	pid := actor.Spawn(props)

	// A map for storing actors
	// hash := make(map[int]messageActor)

	// Reading from file
	file, err := os.Open("text.csv")
	check(err)

	defer file.Close()

	// reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		pid.Tell(&listRow{row: scanner.Text()})
		sentCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Allt är skickat!")
	fmt.Scanln()
	fmt.Printf("Sent: %v\nReceived: %v", sentCounter, receiveCounter)
	if sentCounter == receiveCounter {
		println("\nAlla skickade paket togs emot!")
	} else {
		println("\nNågonting gick fel >:(")
	}

}
