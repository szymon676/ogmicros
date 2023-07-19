package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/ddliu/go-httpclient"
	msg "github.com/szymon676/ogmicros/protos"
)

var (
	msgChan   = make(chan string)
	stockName = flag.String("stock", "TSLA", "stockname for getting recommendations")
)

func fetchData() error {
	res, err := httpclient.Get("http://127.0.0.1:8000/recommendation/" + *stockName)
	if err != nil {
		return err
	}
	bodyString, err := res.ToString()
	if err != nil {
		return err
	}

	msgChan <- bodyString

	return nil
}

func main() {
	flag.Parse()
	e := actor.NewEngine()
	r := remote.New(e, remote.Config{ListenAddr: "127.0.0.1:3000"})
	e.WithRemote(r)
	pid := actor.NewPID("127.0.0.1:4000", "processing-ms")

	go func() {
		var prevData string
		for data := range msgChan {
			fmt.Println("Received data:", data)
			if data != prevData {
				e.Send(pid, &msg.Message{Data: data})
				prevData = data
			}
		}
	}()

	for {
		fetchData()
		time.Sleep(time.Second * 5)
	}
}
