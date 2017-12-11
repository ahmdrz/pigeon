package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tucnak/telebot"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

var (
	pigeon struct {
		port          int
		users         []int
		telegramToken string
		outputFile    *os.File

		bot *telebot.Bot
	}
)

func main() {
	pigeon.users = make([]int, 0)
	pigeon.outputFile = os.Stdout

	err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on parsing arguments , %s", err.Error())
		flag.Usage()

		os.Exit(1)
	}

	logger.SetOutput(pigeon.outputFile)

	logger.Println("Pigeon is starting ...")
	logger.Println("Connecting to telegram ...")

	botSetting := telebot.Settings{
		Token: pigeon.telegramToken,
	}
	pigeon.bot, err = telebot.NewBot(botSetting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Telegram error , %s", err.Error())
		os.Exit(2)
	}

	logger.Println("Connected , Sending `hello` to users ...")

	for _, user := range pigeon.users {
		logger.Printf("Sending to %d", user)
		telegramUser := &telebot.User{ID: user}
		message, err := pigeon.bot.Send(telegramUser, "[pigeon] Hello There ! I'm alive :)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Telegram user error , %s", err.Error())
			os.Exit(3)
		}
		logger.Printf("Message sent to %s %s", message.Chat.FirstName, message.Chat.LastName)
	}

	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterService(new(PigeonService), "")

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", pigeon.port),
		Handler: rpcServer,
	}

	logger.Printf("Listening to %s\n", httpServer.Addr)
	logger.Fatal(httpServer.ListenAndServe())
}

func parseArgs() error {
	flag.IntVar(&pigeon.port, "port", 8765, "Listen port")
	users := flag.String("users", "", "Telegram ID of users. separate it by ','")
	flag.StringVar(&pigeon.telegramToken, "token", "", "Telegram bot token")

	outputPtr := flag.String("output", "", "Output of pigeon.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nNotice: To find your telegram id, you can use @userinfobot\n")
		fmt.Fprintf(os.Stderr, "Example: pigeon-server -port=8595 -token='<BOT-TOKEN>' -users='<USER1>,<USER2>'\n")
	}

	if *users == "" {
		return fmt.Errorf("Please fill users argument \n")
	}

	if *outputPtr != "" {
		var err error
		pigeon.outputFile, err = os.OpenFile(*outputPtr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	}

	if pigeon.telegramToken == "" {
		return fmt.Errorf("Please fill token argument \n")
	}

	tmpList := strings.Split(*users, ",")
	for _, user := range tmpList {
		tmpUser, err := strconv.Atoi(user)
		if err != nil {
			return err
		}
		pigeon.users = append(pigeon.users, tmpUser)
	}

	return nil
}
