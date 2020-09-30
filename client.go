package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
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
	Port		string
	Version     bool
	Rules       []string
}

type Response struct {
	Filename string        `json:"filename"`
	Display  string        `json:"display"`
	Errors   []interface{} `json:"errors"`
}

func (client *Client) Resp(body []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (client *Client) Init() {
	MQConnect, err := amqp.Dial(fmt.Sprint("amqp://", client.Login, ":", client.Password, "@", client.Hostname, ":", client.Port, "/"))
	if err != nil {
		log.Fatal(err)
	}
	client.MQConnect = MQConnect
	ch, err := client.MQConnect.Channel()
	if err != nil {
		log.Fatal(err)
	}
	client.Channel = ch
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
	client.reply_queue = q
	client.Count = 0
	msgs, err := client.Channel.Consume(
		client.reply_queue.Name, "", true,
		false, false, false, nil)
	client.msgs = msgs
	if err != nil {
		log.Print(err)
	}
}

func (client *Client) publish(content []byte) {
	client.Count++
	client.Channel.Publish(
		"",
		"norminette",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			ReplyTo:     client.reply_queue.Name,
			Body:        content,
		})
}

func (client *Client) RequestPreparation(filepath string) ([]byte, error) {
	content, err := ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(File{Name: filepath, Content: content, Rules: client.Rules})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *Client) RequestVersion() ([]byte, error) {
	result, err := json.Marshal(ActiveVersion{Action: "version"})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *Client) SendFile(filepath string) error {
	file, err := client.RequestPreparation(filepath)
	if err != nil {
		return err
	}
	client.publish(file)
	return nil
}

func (client *Client) PrintVersion() {
	for msg := range client.msgs {
		client.Count--
		result, err := client.Resp(msg.Body)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println("Norminette version:")
		fmt.Println(result.Display)
		break
	}
}

func (client *Client) SendVersion() error {
	file, err := client.RequestVersion()
	if err != nil {
		return err
	}
	client.publish(file)
	client.PrintVersion()
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
		if strings.HasPrefix(result.Display, "Unvalid options for -R:") == true {
			fmt.Print(result.Display, "\n")
			os.Exit(0)
		} else {
			fmt.Printf("Norme: %s\n", result.Filename)
			if result.Display != "" {
				fmt.Print(result.Display, "\n")
			}
		}
	}
}
