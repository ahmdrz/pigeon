package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
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
