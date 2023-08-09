package chat

type message struct {
	RoomId string `json:"room_id"`
	Msg    string `json:"msg"`
}

type operation struct {
	roomId string
	clientInfo
}

type Service struct {
	roomChannels map[string]*chatRoom
	open         chan *operation
	close        chan *operation
	delete       chan string
	messages     chan *message
}

func Init() *Service {
	service := &Service{
		roomChannels: make(map[string]*chatRoom),
		open:         make(chan *operation, 100),
		close:        make(chan *operation, 100),
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

func (s *Service) room(roomId string) *chatRoom {
	b, ok := s.roomChannels[roomId]
	if !ok {
		b = makeRoom()
		s.roomChannels[roomId] = b
	}
	return b
}

func (s *Service) register(opened *operation) {
	s.room(opened.roomId).joinRoom(opened.clientInfo)
}

func (s *Service) deregister(closed *operation) {
	s.room(closed.roomId).exitRoom(closed.clientInfo)
	close(closed.stream)
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

func (s *Service) JoinRoom(roomId string, userId string) chan string {
	stream := make(chan string)
	s.open <- &operation{
		roomId: roomId,
		clientInfo: clientInfo{
			userId: userId,
			stream: stream,
		},
	}
	return stream
}

func (s *Service) ExitRoom(roomId string, userId string, stream chan string) {
	s.close <- &operation{
		roomId: roomId,
		clientInfo: clientInfo{
			userId: userId,
			stream: stream,
		},
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

func (s *Service) GetRoomList() []string {
	keys := make([]string, 0, len(s.roomChannels))
	for key := range s.roomChannels {
		keys = append(keys, key)
	}
	return keys
}
