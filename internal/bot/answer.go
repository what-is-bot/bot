package bot

type Answer struct {
	ID     string
	Text   string
	Author string
	Score  int
}

type AnswerProvider interface {
	Ask(Question) ([]Answer, error)
}
