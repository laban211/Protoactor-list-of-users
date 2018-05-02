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

type messageActor struct {
	localCounter int
}

var sentCounter int
var createdKeys int

func (state *messageActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *listRow:
		state.localCounter++
		if state.localCounter == 40 {
			fmt.Printf("\n + %v \n", msg.row)
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
	//pid := actor.Spawn(props)

	// A map for storing actors
	hash := make(map[string]*actor.PID)

	// Reading from file
	file, err := os.Open("text.csv")
	check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		projNum := strings.Split(scanner.Text(), ",")[7]

		//if not exist
		value, ok := hash[projNum]
		if !ok {
			createdKeys++
			value = actor.Spawn(props) //Vill deklarera med new()
			hash[projNum] = value
		}

		value.Tell(&listRow{row: scanner.Text()})

		//pid.Tell(&listRow{row: scanner.Text()})
		sentCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Allt är skickat!")
	fmt.Scanln()
	fmt.Printf("Sent: %v", sentCounter)
	fmt.Printf("\nCreated keys: %v", createdKeys)

}
