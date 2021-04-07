package feedback

type FeedbackWasLeftEvent interface {
	GetFeedback() Feedback
}

type feedbackWasLeftEvent struct {
	feedback Feedback
}

func (s feedbackWasLeftEvent) GetFeedback() Feedback {
	return s.feedback
}

func NewFeedbackWasLeftEvent(f Feedback) FeedbackWasLeftEvent {
	return feedbackWasLeftEvent{
		feedback: f,
	}
}