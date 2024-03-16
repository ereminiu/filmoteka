package controller

type outputWithMessage struct {
	Message string `json:"message"`
}

type outputWithId struct {
	Id      int    `json:"id,omitempty"`
	Message string `json:"message"`
}
