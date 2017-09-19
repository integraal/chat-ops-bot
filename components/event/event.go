package event

import (
	"github.com/integraal/ics-golang"
	"github.com/integraal/chat-ops-bot/components/user"
	"errors"
	"time"
	"strings"
)

var events map[string]*Event = make(map[string]*Event)

type Event struct {
	ID          string
	ImportedID  string
	Summary     string
	Description string
	Duration    time.Duration
	Start       time.Time
	End         time.Time

	ReminderSent bool
	PollSent     bool

	users        map[int64]user.User
	agreed       map[int64]bool
	disagreed    map[int64]bool
}

func (e *Event) SetAgree(userId int64) {
	e.agreed[userId] = true
	if _, ok := e.disagreed[userId]; ok {
		delete(e.disagreed, userId)
	}
}
func (e *Event) SetDisagree(userId int64) {
	e.disagreed[userId] = true
	if _, ok := e.agreed[userId]; ok {
		delete(e.agreed, userId)
	}
}

func (e *Event) GetUser(chatId int64) (*user.User, error) {
	if u, ok := e.users[chatId]; ok {
		return &u, nil
	}
	return nil, errors.New("User does not exist")
}

func (e *Event) GetUsers() *map[int64]user.User {
	return &e.users
}

func Append(event *Event, user user.User) {
	if e, ok := events[event.ID]; ok {
		events[e.ID].users[int64(user.TelegramId)] = user
	} else {
		event.users[int64(user.TelegramId)] = user
		events[event.ID] = event
	}
	// TODO: deal with event updates. When event is updated it gets new ID and old event doesn't disappear
	// TODO: clean old events from memory
}

func NewEvent(ics *ics.Event) Event {
	evt := Event{
		ID:          ics.GetID(),
		ImportedID:  strings.Replace(ics.GetImportedID(), "-", "", -1),
		Summary:     ics.GetSummary(),
		Description: ics.GetDescription(),
		Duration:    ics.GetEnd().Sub(ics.GetStart()),
		Start:       ics.GetStart(),
		End:         ics.GetEnd(),

		users:     make(map[int64]user.User),
		agreed:    make(map[int64]bool),
		disagreed: make(map[int64]bool),
	}
	return evt
}

func GetAll() *map[string]*Event {
	return &events
}

func Get(id string) (*Event, error) {
	if e, ok := events[id]; ok {
		return e, nil
	}
	return nil, errors.New("Event does not exist")
}

func (e *Event) GetAgreedCount() int {
	return len(e.agreed)
}

func (e *Event) GetDisagreedCount() int {
	return len(e.disagreed)
}
