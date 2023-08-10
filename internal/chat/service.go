package chat

import (
	"encoding/json"
	"github.com/google/uuid"
)

type event struct {
	roomId string
	client
}

type message struct {
	roomId string
	Msg    string
}

type Service struct {
	roomChannels map[string]*room
	open         chan *event
	close        chan *event
	delete       chan string
	messages     chan *message

	roomList []Room
}

func Init() *Service {
	service := &Service{
		roomChannels: make(map[string]*room),
		open:         make(chan *event, 100),
		close:        make(chan *event, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *message, 100),
	}
	service.initDefaultRooms()

	go service.run()
	return service
}

func (s *Service) initDefaultRooms() {
	s.roomList = []Room{
		{
			Name: "Emoji Room",
			UUID: uuid.New().String(),
			Desc: "This is an emoji chat room.\nWhen you enter text, it will be changed into appropriate emojis.",
		},
		{
			Name: "Test Room 1",
			UUID: uuid.New().String(),
			Desc: "This is just test room for chatting",
		},
		{
			Name: "Test Room 2",
			UUID: uuid.New().String(),
			Desc: "This is just test room for chatting",
		},
		{
			Name: "Test Room 3",
			UUID: uuid.New().String(),
			Desc: "This is just test room for chatting",
		},
		{
			Name: "Test Room 4",
			UUID: uuid.New().String(),
			Desc: "This is just test room for chatting",
		},
		{
			Name: "Test Room 5",
			UUID: uuid.New().String(),
			Desc: "This is just test room for chatting",
		},
		{
			Name: "Test Room 6",
			UUID: uuid.New().String(),
			Desc: "This is just test room for chatting",
		},
	}
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
			s.room(msg.roomId).broadcast(msg.Msg)
		}
	}
}

func (s *Service) room(roomId string) *room {
	r, ok := s.roomChannels[roomId]
	if !ok {
		r = makeRoom()
		s.roomChannels[roomId] = r
	}
	return r
}

func (s *Service) register(opened *event) {
	s.room(opened.roomId).joinRoom(opened.client)
}

func (s *Service) deregister(closed *event) {
	s.room(closed.roomId).exitRoom(closed.client)
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
	s.open <- &event{
		roomId: roomId,
		client: client{
			userId: userId,
			stream: stream,
		},
	}
	return stream
}

func (s *Service) ExitRoom(roomId string, userId string, stream chan string) {
	s.close <- &event{
		roomId: roomId,
		client: client{
			userId: userId,
			stream: stream,
		},
	}
}

func (s *Service) DeleteBroadcast(roomId string) {
	s.delete <- roomId
}

func (s *Service) SendMessage(roomId string, sender string, text string) {
	msg := Message{
		Sender: sender,
		Msg:    text,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return
	}
	newMessage := &message{
		roomId: roomId,
		Msg:    string(jsonData),
	}
	s.messages <- newMessage
}

func (s *Service) GetRoomList() []Room {
	return s.roomList
}

func (s *Service) GetRoomInfo(roomId string) RoomInfo {
	return s.room(roomId).info
}
