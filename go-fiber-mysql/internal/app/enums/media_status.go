package enums

type MediaStatus string

const (
	MediaStatusPending MediaStatus = "pending"
	MediaStatusSuccess MediaStatus = "success"
	MediaStatusFailed  MediaStatus = "failed"
)
