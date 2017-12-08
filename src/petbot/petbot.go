package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
	"petbot/command_interpreter"

	_ "github.com/go-sql-driver/mysql"
	"log"
	"petbot/models"
)

func main() {

	mysqlUser := "petbotadmin"
	mysqlPassword := "notpetbot"
	mysqlDatabaseName := "petbot"
	dataSourceName := mysqlUser + ":" + mysqlPassword + "@/" + mysqlDatabaseName

	db, err := models.Init("mysql", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if !db.HasPetTable() {
		db.CreatePetTable()
	}

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Printf("Connection counter: %s\n", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message received")
				command_interpreter.InterpretCommand(ev, rtm, db)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}
