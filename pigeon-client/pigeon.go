package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/rpc/json"
)

var (
	pigeon struct {
		port   int
		server string

		title string
		body  string
	}
)

func main() {
	var err error

	err = parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on parsing arguments %s", err.Error())
		flag.PrintDefaults()

		os.Exit(1)
	}

	url := fmt.Sprintf("http://%s:%d", pigeon.server, pigeon.port)
	args := &RequestArgs{Title: pigeon.title, Body: pigeon.body}

	logger.Printf("Sending message, %s\n", args.Hash())

	message, err := json.EncodeClientRequest("PigeonService.Notify", args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on coding request %s", err.Error())
		os.Exit(2)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on makeing request %s", err.Error())
		os.Exit(3)
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on sending request %s", err.Error())
		os.Exit(3)
	}
	defer resp.Body.Close()

	var result RequestReply
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not decode response %s", err.Error())
		os.Exit(4)
	}

	if result.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Invalid status code from server %d", result.StatusCode)
		os.Exit(5)
	}

	logger.Println("Message has been sent successfully")
}

func parseArgs() error {
	flag.StringVar(&pigeon.server, "server", "localhost", "pigeon server address")
	flag.StringVar(&pigeon.title, "title", "", "title of message")
	flag.StringVar(&pigeon.body, "body", "", "body|content of message")
	flag.IntVar(&pigeon.port, "port", 8765, "pigeon server port")

	flag.Parse()

	if len(pigeon.title)*len(pigeon.body) == 0 {
		return fmt.Errorf("title or body is empty\n")
	}
	return nil
}
