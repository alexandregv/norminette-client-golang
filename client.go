package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Client struct {
	MQConnect   *amqp.Connection
	Channel     *amqp.Channel
	reply_queue amqp.Queue
	msgs        <-chan amqp.Delivery
	Count       int
	Hostname    string
	Login       string
	Password    string
}

type Response struct {
	Filename string        `json:"filename"`
	Display  string        `json:"display"`
	Errors   []interface{} `json:"errors"`
}

func (c *Client) Resp(body []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Client) Init() {
	c.Hostname = "norminette.21-school.ru"
	c.Login = "guest"
	c.Password = "guest"
	MQConnect, err := amqp.Dial(fmt.Sprint("amqp://", c.Login, ":", c.Password, "@", c.Hostname, ":5672", "/"))
	if err != nil {
		log.Fatal(err)
	}
	c.MQConnect = MQConnect
	ch, err := c.MQConnect.Channel()
	if err != nil {
		log.Fatal(err)
	}
	c.Channel = ch
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		log.Fatal(err)
	}
	c.reply_queue = q
	c.Count = 0
	msgs, err := c.Channel.Consume(
		c.reply_queue.Name, "", true,
		false, false, false, nil)
	c.msgs = msgs
	if err != nil {
		log.Print(err)
	}
}

func (c *Client) publish(contant []byte) {
	c.Count++
	c.Channel.Publish(
		"",
		"norminette",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			ReplyTo:     c.reply_queue.Name,
			Body:        contant,
		})
}

func (c *Client) SendFile(filepath string) error {
	file, err := RequestPreparation(filepath)
	if err != nil {
		return err
	}
	c.publish(file)
	return nil
}

func (client *Client) PrintResult() {
	for msg := range client.msgs {
		client.Count--
		result, err := client.Resp(msg.Body)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("Norme: %s\n", result.Filename)
		fmt.Print(result.Display, "\n")
	}
}
