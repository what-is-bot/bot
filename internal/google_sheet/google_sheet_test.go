// +build integration

package googlesheet

import (
	"testing"

	"github.com/what-is-bot/bot/internal/bot"
)

func TestGoogleSheetAnswerProvider_Ask(t *testing.T) {
	answerProvider := NewAnswerProvider(defaultSheetsService())

	question := bot.Question{Term: "teste"}
	answers, err := answerProvider.Ask(question)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	for _, answer := range answers {
		t.Logf("answer: %+v", answer)
	}
}

func TestGoogleSheetFeedbackProvider_Upvote(t *testing.T) {
	feedbackProvider := NewFeedbackProvider(defaultSheetsService())

	err := feedbackProvider.Upvote(bot.Answer{ID: "2"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGoogleSheetFeedbackProvider_Downvote(t *testing.T) {
	feedbackProvider := NewFeedbackProvider(defaultSheetsService())

	err := feedbackProvider.Downvote(bot.Answer{ID: "2"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
