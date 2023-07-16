package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	msg "github.com/szymon676/ogmicros/protos"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type server struct{}

func newServer() actor.Receiver {
	return &server{}
}

type data struct {
	time           string
	stockName      string
	recommendation string
}

type StockData struct {
	Data struct {
		StockName      string `json:"stockname"`
		Recommendation string `json:"recommendation"`
	} `json:"data"`
}

func newDataFromMsg(msg *msg.Message) *data {
	var stockData StockData
	if err := json.Unmarshal([]byte(msg.Data), &stockData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return &data{
		time:           time.Now().String(),
		stockName:      stockData.Data.StockName,
		recommendation: stockData.Data.Recommendation,
	}
}

func (f *server) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		fmt.Println("server has started")
	case *actor.PID:
		fmt.Println("server has received:", m)
	case *msg.Message:
		// sendSms(m.Data)
		data := newDataFromMsg(m)
		if data != nil {
			saveToSheet(data)
		}
		fmt.Println("received message:", m)
	}
}

func main() {
	e := actor.NewEngine()
	r := remote.New(e, remote.Config{ListenAddr: "127.0.0.1:4000"})
	e.WithRemote(r)

	e.Spawn(newServer, "processing-ms")
	select {}
}
func saveToSheet(data *data) {
	f, err := os.Create("data.csv")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	records := [][]string{
		{data.time, data.stockName, data.recommendation},
	}

	if err := w.WriteAll(records); err != nil {
		log.Fatal("Error writing to CSV:", err)
	}
}

func sendSms(msg string) {
	accountSid := os.Getenv("SID")
	authToken := os.Getenv("TOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(os.Getenv("TO"))
	params.SetFrom(os.Getenv("FROM"))
	params.SetBody(msg)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}
