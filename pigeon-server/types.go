package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/tucnak/telebot"
)

type RequestArgs struct {
	Title string
	Body  string
}

func (r *RequestArgs) Hash() string {
	hasher := md5.New()
	hasher.Write([]byte(r.Title + r.Body))
	return fmt.Sprintf("%d:%s", time.Now().Unix(), hex.EncodeToString(hasher.Sum(nil)))
}

type RequestReply struct {
	StatusCode int
}

type PigeonService struct{}

func (h *PigeonService) Notify(r *http.Request, args *RequestArgs, reply *RequestReply) error {
	argsHash := args.Hash()

	for _, user := range pigeon.users {
		logger.Printf("Sending message to %d , %s", user, argsHash)
		telegramUser := &telebot.User{ID: user}
		_, err := pigeon.bot.Send(telegramUser, fmt.Sprintf("[pigeon] %s\n%s", args.Title, args.Body))
		if err != nil {
			logger.Printf("Error on sending message %s -> %d", argsHash, user)
			continue
		}
		logger.Printf("Message %s sent to %d", argsHash, user)
	}

	reply.StatusCode = http.StatusOK
	return nil
}
