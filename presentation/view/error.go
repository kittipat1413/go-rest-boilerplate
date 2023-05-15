package view

type Error struct {
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	HttpCode int         `json:"-"`
}
