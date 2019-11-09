package slack

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"github.com/what-is-bot/bot/internal"
)

func Start(c internal.Controller) {
	slackToken := viper.GetString("slack_token")
	fmt.Printf("Slack token %s \n", slackToken)

	api := slack.New(slackToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "GQE33K0R4"))

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "GQE33K0R4"))

		case *slack.ReactionAddedEvent:
			fmt.Println("Infos:", ev.Item)
			fmt.Println("Reaction: ", ev.Reaction)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

type channelWriter struct {
	*bytes.Buffer
	channel string
}
