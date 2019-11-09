package internal

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Question struct {
	Term string
}

var NoQuestion = Question{Term: string(rune(21))} // 21 is the NAK ASCII character

type Answer struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Score  int    `json:"score"`
}

type AnswerProvider interface {
	Ask(Question) ([]Answer, error)
}

type QuestionMatcher interface {
	Match(msg string) (Question, error)
}

type basicQuestionMatcher struct {
	greetings []string
	t         transform.Transformer
}

func (b *basicQuestionMatcher) Match(msg string) (Question, error) {
	if msg == "" {
		return NoQuestion, nil
	}
	transformedMsg, _, err := transform.String(b.t, msg)
	if err != nil {
		return NoQuestion, err
	}
	for _, greet := range b.greetings {
		if index := strings.Index(transformedMsg, greet); index > 0 {
			termStartIndex := index + len(greet) + 1
			if termStartIndex > len(transformedMsg) {
				return NoQuestion, nil
			}
			return Question{Term: transformedMsg[termStartIndex:]}, nil
		}
	}
	return NoQuestion, nil
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func NewBasicQuestionMatcher() QuestionMatcher {
	return &basicQuestionMatcher{
		greetings: []string{"que e", "que seria", "what is", "what's"},
		t:         transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC),
	}
}

type FeedbackProvider interface {
	Upvote(Answer) error
	Downvote(Answer) error
}

type Controller interface {
	Answer(msg string, out io.Writer) error
	Upvote(Answer) error
	Downvote(Answer) error
}

type controller struct {
	qMatcher  QuestionMatcher
	aProvider AnswerProvider
	fProvider FeedbackProvider
}

func (c *controller) Answer(msg string, out io.Writer) error {
	jsonOut := json.NewEncoder(out)
	question, err := c.qMatcher.Match(msg)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	if question == NoQuestion {
		log.Println("no question")
		return nil
	}

	answers, err := c.aProvider.Ask(question)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	err = jsonOut.Encode(answers)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	return nil
}

func (c *controller) Upvote(answer Answer) error {
	return c.fProvider.Upvote(answer)
}

func (c *controller) Downvote(answer Answer) error {
	return c.fProvider.Downvote(answer)
}

func NewController(provider AnswerProvider, fProvider FeedbackProvider) Controller {
	return &controller{
		qMatcher:  NewBasicQuestionMatcher(),
		aProvider: provider,
		fProvider: fProvider,
	}
}
