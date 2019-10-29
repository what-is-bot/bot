package bot

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

type Question struct {
	Term string
}

var NoQuestion = Question{Term: string(rune(21))} // 21 is the NAK ASCII character

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
