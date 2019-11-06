package googlesheet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/mitchellh/go-homedir"
	"github.com/what-is-bot/bot/internal/bot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	sheetID = "10IwyKQ4qd8eyeqd-VGI5bvJjy28Wu2Ja6VcUtEogLuE"
)

type feedbackProvider struct {
	srv *sheets.Service
}

func (g *feedbackProvider) Upvote(answer bot.Answer) error {
	rowIndex, _ := strconv.Atoi(answer.ID)
	return writeVote(g.srv, rowIndex, 1)
}

func (g *feedbackProvider) Downvote(answer bot.Answer) error {
	rowIndex, _ := strconv.Atoi(answer.ID)
	return writeVote(g.srv, rowIndex, -1)
}

func writeVote(srv *sheets.Service, rowIndex, delta int) error {
	value, err := readValue(srv, rowIndex)
	if err != nil {
		return err
	}
	score, err := strconv.Atoi(value[3].(string))
	if err != nil {
		return err
	}

	var updateValues []interface{}

	updateValues = append(updateValues, score+delta)

	err = writeValue(srv, rowIndex, updateValues)
	if err != nil {
		return err
	}

	return nil
}

func writeValue(srv *sheets.Service, rowIndex int, values []interface{}) error {
	rng := fmt.Sprintf("answers!D%d", rowIndex)
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         [][]interface{}{values},
	}
	_, err := srv.Spreadsheets.Values.
		Update(sheetID, rng, &valueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return err
	}
	return nil
}

func readValue(srv *sheets.Service, rowIndex int) ([]interface{}, error) {
	resp, err := srv.Spreadsheets.Values.Get(sheetID, fmt.Sprintf("answers!A%d:D%d", rowIndex, rowIndex)).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values[0], nil
}

func NewFeedbackProvider(srv *sheets.Service) bot.FeedbackProvider {
	return &feedbackProvider{srv: srv}
}

type answerProvider struct {
	srv *sheets.Service
}

func (g *answerProvider) Ask(question bot.Question) ([]bot.Answer, error) {
	resp, err := g.srv.Spreadsheets.Values.Get(sheetID, "answers!A2:D").Do()
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Values) == 0 {
		return nil, nil
	}

	answers := filterValues(resp.Values, func(t string) bool { return t == question.Term })

	return answers, nil
}

func filterValues(rows [][]interface{}, predicate func(term string) bool) []bot.Answer {
	var result []bot.Answer
	for idx, row := range rows {
		if predicate(row[0].(string)) {
			score, _ := strconv.Atoi(row[3].(string))
			result = append(result, bot.Answer{
				ID:     strconv.Itoa(idx + 2),
				Text:   row[1].(string),
				Author: row[2].(string),
				Score:  score,
			})
		}
	}

	return result
}

func defaultSheetsService() *sheets.Service {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	b, err := ioutil.ReadFile(path.Join(home, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}

func NewAnswerProvider(srv *sheets.Service) bot.AnswerProvider {
	return &answerProvider{srv}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	tokFile := "token.json"
	tok, err := tokenFromFile(path.Join(home, tokFile))
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
