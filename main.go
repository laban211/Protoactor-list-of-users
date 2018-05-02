package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type listRow struct {
	row string
}

type messageActor struct {
	localCounter int
}

type askForRows struct {
	test int
}

var sentCounter int
var createdKeys int
var receivedCounter int

func (state *messageActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *listRow:
		state.localCounter++
		if state.localCounter > 100 {
			fmt.Printf("\n%v\n", msg.row)
		}
	case *askForRows:
		// fmt.Println("hej")
		context.Respond(state.localCounter)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	println("Skickar...\n")

	// Create an actor
	props := actor.FromProducer(func() actor.Actor { return &messageActor{} })

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
			value = actor.Spawn(props)
			hash[projNum] = value
		}

		value.Tell(&listRow{row: scanner.Text()})

		sentCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for _, value := range hash {
		result, _ := value.RequestFuture(&askForRows{test: 1}, 30*time.Second).Result() // await result
		receivedCounter += result.(int)
	}

	fmt.Println("Allt är skickat!")
	// fmt.Scanln()
	fmt.Printf("Sent: %v\nCreated keys: %v\nReceived: %v",
		sentCounter, createdKeys, receivedCounter)
}
