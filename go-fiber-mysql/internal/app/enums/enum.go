package enums

type Enum interface {
	IsValid() bool
	GetValues() []string
}
