package chat

type message struct {
	RoomId string `json:"room_id"`
	Msg    string `json:"msg"`
}

type listener struct {
	RoomId string
	Chan   chan string
}

type Service struct {
	roomChannels map[string]IChatRoom
	open         chan *listener
	close        chan *listener
	delete       chan string
	messages     chan *message
}

func Init() *Service {
	service := &Service{
		roomChannels: make(map[string]IChatRoom),
		open:         make(chan *listener, 100),
		close:        make(chan *listener, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *message, 100),
	}

	go service.run()
	return service
}

func (s *Service) run() {
	for {
		select {
		case opened := <-s.open:
			s.register(opened)
		case closed := <-s.close:
			s.deregister(closed)
		case roomId := <-s.delete:
			s.deleteRoom(roomId)
		case msg := <-s.messages:
			s.room(msg.RoomId).sendMessage(msg.Msg)
		}
	}
}

func (s *Service) room(roomId string) IChatRoom {
	b, ok := s.roomChannels[roomId]
	if !ok {
		b = makeRoom()
		s.roomChannels[roomId] = b
	}
	return b
}

func (s *Service) register(opened *listener) {
	s.room(opened.RoomId).joinRoom(opened.Chan)
}

func (s *Service) deregister(closed *listener) {
	s.room(closed.RoomId).exitRoom(closed.Chan)
	close(closed.Chan)
}

func (s *Service) deleteRoom(roomId string) {
	b, ok := s.roomChannels[roomId]
	if ok {
		err := b.closeRoom()
		if err != nil {
			return
		}
		delete(s.roomChannels, roomId)
	}
}

func (s *Service) JoinRoom(roomId string) chan string {
	newListener := make(chan string)
	s.open <- &listener{
		RoomId: roomId,
		Chan:   newListener,
	}
	return newListener
}

func (s *Service) ExitRoom(roomId string, closed chan string) {
	s.close <- &listener{
		RoomId: roomId,
		Chan:   closed,
	}
}

func (s *Service) DeleteBroadcast(roomId string) {
	s.delete <- roomId
}

func (s *Service) SendMessage(roomId string, text string) {
	msg := &message{
		RoomId: roomId,
		Msg:    text,
	}
	s.messages <- msg
}
