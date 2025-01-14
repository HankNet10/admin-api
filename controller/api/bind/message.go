package bind

type Message struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}
type ErrorMessage struct {
	Message string `json:"message"`
}
