package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// type listRow struct {
// 	week    int
// 	user    string
// 	mon     float32
// 	tue     float32
// 	wed     float32
// 	thu     float32
// 	fri     float32
// 	projNum int
// }

type listRow struct {
	row string
}

type messageActor struct{}

var receiveCounter int
var sentCounter int

func (state *messageActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	// case *listRow:
	// 	fmt.Printf("%v \t %v \t %v \t %v \t %v \t %v \t %v \t %v \n",
	// 		msg.week, msg.user, msg.mon, msg.tue, msg.wed, msg.thu,
	// 		msg.fri, msg.projNum)
	case *listRow:
		receiveCounter++
		if receiveCounter%50000 == 0 {
			fmt.Println(receiveCounter, msg.row)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	props := actor.FromProducer(func() actor.Actor { return &messageActor{} })
	pid := actor.Spawn(props)

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

	// for {
	// 	line, _, err := reader.ReadLine()

	// 	if err == io.EOF {
	// 		break
	// 	}

	// 	pid.Tell(&listRow2{row: string(line)})
	// }

	//fmt.Println(dat) // Print the content as 'bytes'
	// str := string(file)
	// fmt.Println(str) // Print the content as 'string'

	// pid.Tell(&listRow{
	// 	week:    3,
	// 	user:    "Daniel",
	// 	mon:     4,
	// 	tue:     3,
	// 	wed:     2,
	// 	thu:     8,
	// 	fri:     1,
	// 	projNum: 1234567890,
	// })
	// console.ReadLine()
}
