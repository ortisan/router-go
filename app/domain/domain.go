package domain

type Error struct {
	Message    string `json:"message,omitempty"`
	Cause      string `json:"cause,omitempty"`
	StackTrace string `json:"stacktrace,omitempty"`
}
