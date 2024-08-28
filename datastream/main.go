package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// wsDepthHandler := func(event *binance.WsDepthEvent) {
	// 	fmt.Println(event.Bids)
	// }
	errHandler := func(err error) {
		fmt.Println(err)
	}
	// doneC, stopC, err := binance.WsDepthServe("BTCUSDT", wsDepthHandler, errHandler)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // use stopC to exit
	// go func() {
	// 	time.Sleep(5 * time.Minute)
	// 	stopC <- struct{}{}
	// }()
	// // remove this if you do not want to be blocked here
	// <-doneC
	symbol := "BTCUSDT"
	cfg := GetRabbitConfig()
	conn := cfg.CreateConnection()
	channel, err := conn.Channel()
	if err != nil {
		log.Panicln(err)
	}

	sendStream, err := CreateStream(context.Background(), channel, symbol)
	if err != nil {
		log.Panicln(err)
	}

	wsAggTradeHandler := func(event *binance.WsAggTradeEvent) {
		log.Println(event)
		err := sendStream(event)
		if err != nil {
			log.Panicln(err)
		}

	}

	doneC, _, err := binance.WsAggTradeServe(symbol, wsAggTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}

type PublishFunction func(data any) error

func CreateStream(ctx context.Context, ch *amqp091.Channel, symbol string) (PublishFunction, error) {
	err := ch.ExchangeDeclare(
		symbol,
		"fanout", // type
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		true,     // no-wait
		nil,      // arguments
	)

	return func(data any) error {
		body, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = ch.PublishWithContext(ctx,
			symbol, // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        []byte(body),
			})

		return err

	}, err
}
