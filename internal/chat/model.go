package chat

type Room struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	Desc string `json:"desc"`
}

type Message struct {
	Sender string `json:"sender"`
	Msg    string `json:"msg"`
}

type RoomInfo struct {
	Users    map[string]bool `json:"users"`
	Messages []string        `json:"messages"`
}
