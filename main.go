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

type askForRows struct {
	test int
}

var sentCounter int

func (state *mainRouter) Receive(context actor.Context) {
	projects := actor.FromProducer(newProjectActor)
	users := actor.FromProducer(newUserActor)
	children := make(map[string]*actor.PID)

	foundProjectManager, ok := children["projects"]
	if !ok {
		foundProjectManager = actor.Spawn(projects)
		children["projects"] = foundProjectManager
	}

	foundUserManager, ok := children["users"]
	if !ok {
		foundUserManager = actor.Spawn(users)
		children["users"] = foundUserManager
	}

	switch msg := context.Message().(type) {
	case *listRow:
		children["projects"].Tell(msg)
		children["users"].Tell(msg)
	default:
		fmt.Print("Something went wrong >:((")
	}
}

func (state *projectManager) Receive(context actor.Context) {
	props := actor.FromProducer(func() actor.Actor { return &projectActor{} })
	actors := make(map[string]*actor.PID)
	switch msg := context.Message().(type) {
	case *listRow:
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
	println("HI")
}

func (state *userManager) Receive(context actor.Context) {
	props := actor.FromProducer(func() actor.Actor { return &userActor{} })
	users := make(map[string]*actor.PID)
	switch msg := context.Message().(type) {
	case *listRow:
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
	println("Ho")
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
	props := actor.FromProducer(newParentActor)
	pid := actor.Spawn(props)

	// Reading from file
	file, err := os.Open("text.csv")
	check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		pid.Tell(&listRow{row: scanner.Text()})
		sentCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Allt är skickat!")
	// fmt.Scanln()
}
