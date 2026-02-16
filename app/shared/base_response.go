package shared

type BaseResponse struct {
	Success  bool      `json:"success"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Message  string    `json:"message"`
	Data     any       `json:"data,omitempty"`
}

type Metadata struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Pages   int `json:"pages"`
}
