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

type mainRouter struct{}
type projectManager struct{}
type userManager struct{}
type projectActor struct{}
type userActor struct{}

// type messageActor struct {
// 	localCounter int
// }

type askForRows struct {
	test int
}

var sentCounter int
var createdKeys int
var receivedCounter int

// func (state *messageActor) Receive(context actor.Context) {
// 	switch msg := context.Message().(type) {
// 	case *listRow:
// 		state.localCounter++
// 		if state.localCounter > 100 {
// 			fmt.Printf("\n%v\n", msg.row)
// 		}
// 	case *askForRows:
// 		// fmt.Println("hej")
// 		context.Respond(state.localCounter)
// 	}
// }

func (state *mainRouter) Receive(context actor.Context) {
	//Hur deklarerar jag en actor?
	projects := actor.FromProducer(newProjectActor)
	users := actor.FromProducer(newUserActor)
	childActors := make(map[string]*actor.PID)

	foundProjectActor, ok := childActors["projects"]
	if !ok {
		foundProjectActor = actor.Spawn(projects)
		childActors["projects"] = foundProjectActor
	}

	foundUserActor, ok := childActors["users"]
	if !ok {
		foundUserActor = actor.Spawn(users)
		childActors["users"] = foundUserActor
	}

	switch msg := context.Message().(type) {
	case *listRow:
		//tell them what to do
		childActors["projects"].Tell(msg)
		childActors["users"].Tell(msg)
	default:
		fmt.Print("Something went wrong >:((")
	}
}

func (state *projectManager) Receive(context actor.Context) {
	props := actor.FromProducer(func() actor.Actor { return &projectActor{} })
	actors := make(map[string]*actor.PID)
	switch msg := context.Message().(type) {
	case *listRow:
		//Split
		projNum := strings.Split(msg.row, ",")[7]
		foundActor, ok := actors[projNum]
		if !ok {
			foundActor = actor.Spawn(props)
			actors[projNum] = foundActor
		}
		foundActor.Tell(msg)
	}
}

func (state *projectActor) Receive(context actor.Context) {

}

func (state *userManager) Receive(context actor.Context) {
	props := actor.FromProducer(func() actor.Actor { return &userActor{} })
	users := make(map[string]*actor.PID)
	switch msg := context.Message().(type) {
	case *listRow:
		//Split
		userID := strings.Split(msg.row, ",")[1]
		foundUser, ok := users[userID]
		if !ok {
			foundUser = actor.Spawn(props)
			users[userID] = foundUser
		}
		foundUser.Tell(msg)
	}
}

func (state *userActor) Receive(context actor.Context) {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func newParentActor() actor.Actor {
	return &mainRouter{}
}

func newProjectActor() actor.Actor {
	return &projectManager{}
}

func newUserActor() actor.Actor {
	return &userManager{}
}

func main() {
	println("Skickar...\n")

	// Create an actor
	// props := actor.FromProducer(func() actor.Actor { return &messageActor{} })
	props := actor.FromProducer(newParentActor)
	pid := actor.Spawn(props)

	// A map for storing actors
	// hash := make(map[string]*actor.PID)

	// Reading from file
	file, err := os.Open("text.csv")
	check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// projNum := strings.Split(scanner.Text(), ",")[7]

		//if not exist
		// foundActor, ok := hash[projNum]
		// if !ok {
		// 	createdKeys++
		// 	foundActor = actor.Spawn(props)
		// 	hash[projNum] = foundActor
		// }

		// foundActor.Tell(&listRow{row: scanner.Text()})
		pid.Tell(&listRow{row: scanner.Text()})
		sentCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// for _, value := range hash {
	// 	result, _ := value.RequestFuture(&askForRows{test: 1}, 30*time.Second).Result() // await result
	// 	receivedCounter += result.(int)
	// }

	fmt.Println("Allt är skickat!")
	// fmt.Scanln()
	// fmt.Printf("Sent: %v\nCreated keys: %v\nReceived: %v",
	// 	sentCounter, createdKeys, receivedCounter)
}
