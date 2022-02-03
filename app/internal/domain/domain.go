package domain

type Error struct {
	TraceId    string `json:"trace_id,omitempty"`
	Message    string `json:"message,omitempty"`
	Cause      string `json:"cause,omitempty"`
	StackTrace string `json:"stacktrace,omitempty"`
}
