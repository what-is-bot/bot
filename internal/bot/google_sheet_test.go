// +build integration

package bot

import "testing"

func TestGoogleSheetAnswerProvider_Ask(t *testing.T) {
	answerProvider := NewGoogleSheetAnswerProvider(defaultSheetsService())

	question := Question{Term: "teste"}
	answers, err := answerProvider.Ask(question)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	for _, answer := range answers {
		t.Logf("answer: %+v", answer)
	}
}

func TestGoogleSheetFeedbackProvider_Upvote(t *testing.T) {
	feedbackProvider := NewGoogleSheetFeedbackProvider(defaultSheetsService())

	err := feedbackProvider.Upvote(Answer{ID: "2"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGoogleSheetFeedbackProvider_Downvote(t *testing.T) {
	feedbackProvider := NewGoogleSheetFeedbackProvider(defaultSheetsService())

	err := feedbackProvider.Downvote(Answer{ID: "2"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
