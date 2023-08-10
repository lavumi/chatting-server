package chat

import "log"

type room struct {
	message chan string
	open    chan client
	close   chan client
	streams map[chan string]bool
	info    RoomInfo
}

type client struct {
	userId string
	stream chan string
}

func makeRoom() *room {
	r := &room{
		message: make(chan string),
		open:    make(chan client),
		close:   make(chan client),
		streams: make(map[chan string]bool),

		info: RoomInfo{
			Users:    make(map[string]bool),
			Messages: make([]string, 0),
		},
	}

	go r.listen()

	return r
}

func (r *room) listen() {
	for {
		select {
		case c := <-r.open:
			r.streams[c.stream] = true
			r.info.Users[c.userId] = true
		case c := <-r.close:
			delete(r.streams, c.stream)
			delete(r.info.Users, c.userId)
		case eventMsg := <-r.message:
			for clientMessageChan := range r.streams {
				clientMessageChan <- eventMsg
			}
			r.info.Messages = append(r.info.Messages, eventMsg)
			log.Printf("message size : %d", len(r.info.Messages))
		}
	}
}

func (r *room) joinRoom(c client) {
	r.open <- c
}

func (r *room) exitRoom(c client) {
	r.close <- c
}

func (r *room) broadcast(msg string) {
	r.message <- msg
}

func (r *room) closeRoom() error {
	close(r.open)
	close(r.close)
	return nil
}
