package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{}

type FileLogger struct {
	filename string
}

type RemoteLogger struct {
	url string
}

func (c ConsoleLogger) Log(message string) {
	fmt.Println(message)
}
func (f FileLogger) Log(message string) {
	file, err := os.OpenFile(f.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file")
	}
	defer file.Close()
	file.WriteString(message)

}

func (r RemoteLogger) Log(message string) {
	req, err := http.NewRequest("POST", r.url, bytes.NewBuffer([]byte(message)))
	if err != nil {
		fmt.Println("Error creating request")
	}
	req.Header.Set("Accept", "application/json")
	fmt.Println("[Remote] Отправка на сервер", r.url, "->", message)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
	}
	defer resp.Body.Close()
	fmt.Println("[Remote] Ответ сервера", resp.Status)

}

func main() {

	console := ConsoleLogger{}
	file := FileLogger{"log.txt"}
	remote := RemoteLogger{"http://localhost:8080"}

	loggers := []Logger{console, file, remote}

	for _, logger := range loggers {
		logger.Log("Hello, World!")
	}

}
