package main

import (
	"fmt"
	"time"
	"bufio"
	"os"
	"strings"
	"math/rand"
	"sync"
)

const (
    MAX_HALL_SIZE = 5
	MAX_WINDOWS = 3
)

type Customer struct {
	number	int
	name	string
}

var queue = make(chan Customer, MAX_HALL_SIZE)
var wg sync.WaitGroup
var counter = 0

func generateWorkerTime() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(10)
}

func serviceWindow(window int) {
	for {
		customer, ok := <-queue
		if !ok {
			fmt.Printf("Window-%d out of service.\n", window)
			break
		}

		fmt.Printf("Window-%d called number %d\n", window, customer.number)
		serviceTime := generateWorkerTime()
		time.Sleep(time.Second * time.Duration(serviceTime))

		fmt.Printf("Window-%d service customer %s, number: %d, use time: %d\n",
					window, customer.name, customer.number, serviceTime)
	}
	wg.Done()
}

func callNumber(name string) {
	counter += 1
	before := len(queue)
	select {
	case queue <- Customer{number: counter, name: name}:
		fmt.Printf("Customer: %s called number: %d, before: %d\n", name, counter, before)
	default:
		fmt.Println("Hall is full! Please call number later!")
	}
}

func kiosks() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		customerName, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("There ware errors reading, exiting program.")
			return
		}
		customerName = strings.Replace(customerName, "\n", "", -1)
		if len(customerName) == 0 {
			break
		}
		callNumber(customerName)
	}
	close(queue)
	wg.Done()
}

func main() {
	wg.Add(MAX_WINDOWS + 1)

	callNumber("Anna")
	callNumber("Judy")
	callNumber("Mark")

	for i := 0; i < MAX_WINDOWS; i++ {
		go serviceWindow(i+1)
	}

	kiosks()

	wg.Wait()
}

