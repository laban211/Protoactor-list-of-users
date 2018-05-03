package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type mainRouter struct{}
type projectManager struct{}
type userManager struct{}
type projectActor struct{}
type userActor struct {
	userWorktime float64
}
type listRow struct {
	row string
}
type askForUserWorkTime struct {
	User string
}

var children = make(map[string]*actor.PID)
var actors = make(map[string]*actor.PID)
var users = make(map[string]*actor.PID)
var sentCounter int

func (state *mainRouter) Receive(context actor.Context) {
	projects := actor.FromProducer(
		func() actor.Actor { return &projectManager{} })
	users := actor.FromProducer(
		func() actor.Actor { return &userManager{} })

	if foundProjectManager, ok := children["projects"]; !ok {
		foundProjectManager = actor.Spawn(projects)
		children["projects"] = foundProjectManager
	}

	if foundUserManager, ok := children["users"]; !ok {
		foundUserManager = actor.Spawn(users)
		children["users"] = foundUserManager
	}

	switch msg := context.Message().(type) {
	case *listRow:
		children["projects"].Tell(msg)
		children["users"].Tell(msg)

	case *askForUserWorkTime:
		result, _ := children["users"].RequestFuture(
			&askForUserWorkTime{User: msg.User}, 30*time.Second).Result()
		context.Respond(result)
	}
}

func (state *projectManager) Receive(context actor.Context) {
	props := actor.FromProducer(
		func() actor.Actor { return &projectActor{} })
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
	// println("HI")
}

func (state *userManager) Receive(context actor.Context) {
	props := actor.FromProducer(
		func() actor.Actor { return &userActor{} })
	switch msg := context.Message().(type) {
	case *listRow:
		userID := strings.Split(msg.row, ",")[1]
		foundUser, ok := users[userID]
		if !ok {
			foundUser = actor.Spawn(props)
			users[userID] = foundUser
		}
		foundUser.Tell(msg)
	case *askForUserWorkTime:
		fmt.Println("usermanager: ", msg.User)
		fmt.Println("till vem: ", users[msg.User])
		result, _ := users[msg.User].RequestFuture(
			&askForUserWorkTime{User: msg.User}, 30*time.Second).Result()
		context.Respond(result)
	}
}

func (state *userActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *listRow:
		var userWorkTimeList [5]float64
		for i := 0; i < 5; i++ {
			userWorkTimeList[i], _ = strconv.ParseFloat(
				strings.Split(msg.row, ",")[i+2], 64)
			state.userWorktime += userWorkTimeList[i]
		}
	case *askForUserWorkTime:
		context.Respond(state.userWorktime)
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
	props := actor.FromProducer(
		func() actor.Actor { return &mainRouter{} })
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

	result, _ := pid.RequestFuture(
		&askForUserWorkTime{User: "user83"}, 30*time.Second).Result()
	fmt.Println("User83 har jobbat: ", result)

	fmt.Println("Allt är skickat!")
}
