package chat

type chatRoom struct {
	message  chan string
	open     chan clientInfo
	close    chan clientInfo
	streams  map[chan string]bool
	users    map[string]bool
	messages []string
}

type clientInfo struct {
	userId string
	stream chan string
}

func makeRoom() *chatRoom {
	room := &chatRoom{
		message:  make(chan string),
		open:     make(chan clientInfo),
		close:    make(chan clientInfo),
		streams:  make(map[chan string]bool),
		users:    make(map[string]bool),
		messages: make([]string, 0),
	}

	go room.listen()

	return room
}

func (room *chatRoom) listen() {
	for {
		select {
		case client := <-room.open:
			room.streams[client.stream] = true
			room.users[client.userId] = true
		case client := <-room.close:
			delete(room.streams, client.stream)
			delete(room.users, client.userId)
		case eventMsg := <-room.message:
			for clientMessageChan := range room.streams {
				clientMessageChan <- eventMsg
			}
			room.messages = append(room.messages, eventMsg)
		}
	}
}

func (room *chatRoom) joinRoom(client clientInfo) {
	room.open <- client
}

func (room *chatRoom) exitRoom(client clientInfo) {
	room.close <- client
}

func (room *chatRoom) sendMessage(msg string) {
	room.message <- msg
}

func (room *chatRoom) closeRoom() error {
	close(room.open)
	close(room.close)
	return nil
}
