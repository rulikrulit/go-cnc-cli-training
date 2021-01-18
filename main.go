package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

// Dest where the message should be sent
type Dest string

const (
	// Email destination for messages
	Email Dest = "email"
	// Phone destination for messages
	Phone Dest = "phone"
	// All destinations for messages
	All Dest = "all"
)

// Message to be sent
type Message struct {
	ID    string `json:"_id"`
	By    Dest   `json:"by"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Body  string `json:"body"`
}

// do your processing here don't forget to use waitGroup
func emailWorker(message Message, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("sent email message " + message.ID)
}

// do your processing here don't forget to use waitGroup
func phoneWorker(message Message, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(200 * time.Millisecond)
	fmt.Println("sent text message " + message.ID)
}

func main() {
	// first receive -f flag from the CLI
	var msgs []Message
	var wg sync.WaitGroup

	filePath := flag.String("f", "messages.json", "file path for messages")
	flag.Parse()

	// then read all the messages from the file and parse JSON
	f, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Println(err)
		return
	}

	if err = json.Unmarshal(f, &msgs); err != nil {
		log.Println(err)
		return
	}

	// create your wait group and channels
	wg.Add(1)
	phoneChan := make(chan Message)
	emailChan := make(chan Message)

	// spin up your workers (don't forget to increment
	// waitGroup before spinning up every worker)
	go func() {
		log.Println("start getting text msgs")
		for msg := range phoneChan {
			phoneWorker(msg, &wg)
		}
	}()

	go func() {
		log.Println("start getting email msgs")
		for msg := range emailChan {
			emailWorker(msg, &wg)
		}
	}()

	// send all you messages (don't forget to check how the message
	// needs to be delivered before sending it in the channel)
	go func() {
		log.Println("start sending msgs")
		log.Println(len(msgs))
		defer wg.Done()
		for n, msg := range msgs {
			if msg.By == Phone {
				log.Println("sending text", n)
				wg.Add(1)
				phoneChan <- msg
			}
			if msg.By == Email {
				log.Println("sending email", n)
				wg.Add(1)
				emailChan <- msg
			}
			if msg.By == All {
				wg.Add(2)
				emailChan <- msg
				phoneChan <- msg
			}
		}
	}()

	wg.Wait()
	log.Println("closing channels")
	close(phoneChan)
	close(emailChan)
	// close your channels and set waitGroup
	// to wait till all workers is finished
}
