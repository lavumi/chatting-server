package chat

type Message struct {
	RoomId string `json:"room_id"`
	Sender string `json:"sender"`
	Msg    string `json:"msg"`
}

type Listener struct {
	RoomId string
	Chan   chan interface{}
}

type Service struct {
	roomChannels map[string]Room
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
}
