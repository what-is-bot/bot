package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

func Start() {
	slackToken := viper.GetString("slack.token")
	fmt.Printf("Slack token %s \n", slackToken)

	api := slack.New(slackToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	matcher := NewBasicQuestionMatcher()
	answerProvider := googlesheet.NewAnswerProvider()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			question, err := matcher.Match(ev.Msg.Text)
			if err != nil {
				fmt.Println(err)
				break
			}
			if question != NoQuestion {
				ans, err := answerProvider.Ask(question)
				if err != nil {
					log.Fatal(err)
				}
				for _, a := range ans {
					log.Printf("answer: %s, score: %d", a.Text, a.Score)
				}
			}

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
