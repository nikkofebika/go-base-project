package enums

type MediaType string

const (
	MediaTypeActivity MediaType = "activity"
	MediaTypeComment  MediaType = "comment"
	MediaTypeAsset    MediaType = "asset"
	MediaTypeTicket   MediaType = "ticket"
	MediaTypeUser     MediaType = "user"
)
