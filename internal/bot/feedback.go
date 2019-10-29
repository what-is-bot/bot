package bot

type FeedbackProvider interface {
	Upvote(Answer) error
	Downvote(Answer) error
}
