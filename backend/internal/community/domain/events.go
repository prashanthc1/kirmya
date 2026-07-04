package domain

// Domain event type names published by the communities context onto the event bus.
const (
	EventPostCreated  = "CommunityPostCreated"
	EventCommentAdded = "CommunityCommentAdded"
	EventPollVoted    = "CommunityPollVoted"
	EventPostReported = "CommunityPostReported"
	EventPostHidden   = "CommunityPostHidden"
)
