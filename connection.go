package homeassistant

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Connection .
type Connection struct {
	sync.RWMutex
	hostname    string
	port        int
	accessToken string
	secure      bool
	conn        *websocket.Conn
	handlers    map[uint64]interface{}
	lastID      uint64
	Done        chan struct{}
}

// NewConnection .
func NewConnection(hostname string, port int, accessToken string, secure bool) (*Connection, error) {
	ha := new(Connection)
	ha.hostname = hostname
	ha.port = port
	ha.accessToken = accessToken
	ha.secure = secure
	ha.handlers = make(map[uint64]interface{})

	go func() {
		for {
			err := backoff.Retry(ha.connectAndAuthenticate, backoff.NewExponentialBackOff())
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}

			ha.readMessages()

			ha.Lock()
			for id, handler := range ha.handlers {
				go handler.(CloseHandler).HandleClose(ha, id)
				delete(ha.handlers, id)
			}
			ha.conn = nil
			ha.Unlock()
		}
	}()

	return ha, nil
}

func (ha *Connection) connectAndAuthenticate() error {
	ha.Lock()
	defer ha.Unlock()

	var err error

	var u url.URL
	if ha.secure {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Host = ha.hostname
	if ha.port == 0 {
		u.Host = ha.hostname
	} else {
		u.Host = fmt.Sprintf("%s:%d", ha.hostname, ha.port)
	}
	u.Path = "/api/websocket"

	log.Printf("attempting to connect to %s", u.String())

	ha.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("****** %+v", err)
		return errors.Wrap(err, "unable to connect")
	}

	for {
		log.Printf("****** about to read message")
		var message Message
		err := ha.conn.ReadJSON(&message)
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					return errors.Wrap(err, "connection closed")
				default:
					return errors.Wrap(err, "some other connection closed error")
				}
			} else {
				return errors.Wrap(err, "problem reading message")
			}

		}
		log.Printf("%+v", message)
		switch message.Type {
		case "auth_required":
			var message = struct {
				Type        string `json:"type"`
				AccessToken string `json:"access_token"`
			}{
				Type:        "auth",
				AccessToken: ha.accessToken,
			}

			err = ha.conn.WriteJSON(&message)
			if err != nil {
				return errors.Wrap(err, "unable to write")
			}

		case "auth_ok":
			log.Printf("authenticated")
			return nil

		case "auth_invalid":
			ha.conn.Close()
			log.Printf("authentication failed")
			return errors.Errorf("authentication failed")
		}
	}

}

func (ha *Connection) readMessages() {
	for {
		var message Message
		err := ha.conn.ReadJSON(&message)
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					log.Printf("CLOSED!")
					return
				}
			}
			log.Printf("read error: %+v", err)
			return
		}

		switch message.Type {
		case "result":
			ha.handleResult(message)

		case "event":
			ha.handleEvent(message)

		case "pong":
			ha.handlePong(message)
			log.Printf("pong\n")

		default:
			log.Printf("unkown message type '%s'", message.Type)
		}
	}
}

func (ha *Connection) sendMessage(handler interface{}, message Command) (uint64, error) {

	ha.Lock()
	defer ha.Unlock()

	if ha.conn == nil {
		return 0, errors.Errorf("connection is not open")
	}

	ha.lastID++
	message.SetID(ha.lastID)

	if handler != nil {
		ha.handlers[ha.lastID] = handler
	}

	err := ha.conn.WriteJSON(&message)
	if err != nil {
		delete(ha.handlers, ha.lastID)
		return 0, errors.Wrap(err, "unable to write message")
	}

	return ha.lastID, nil
}

// RemoveHandler .
func (ha *Connection) RemoveHandler(id uint64) {
	ha.Lock()
	defer ha.Unlock()
	delete(ha.handlers, id)
}

// Close cleanly closes the connection by sending a close message and then
// waiting (with timeout) for the server to close the connection.
func (ha *Connection) Close() {
	ha.Lock()
	err := ha.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"",
		),
	)
	ha.Unlock()
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-ha.Done:
	case <-time.After(time.Second):
	}
}

func (ha *Connection) handlePong(message Message) {
	ha.RLock()
	defer ha.RUnlock()

	if handler, ok := ha.handlers[message.ID]; ok {
		go handler.(PongHandler).HandlePong(ha, message.ID)
	} else {
		log.Printf("no handler registered for %d", message.ID)
	}
}

func (ha *Connection) handleResult(message Message) {
	ha.RLock()
	defer ha.RUnlock()

	if handler, ok := ha.handlers[message.ID]; ok {
		go handler.(ResultHandler).HandleResult(ha, message.ID, message.Success, message.Result)
	} else {
		log.Printf("no handler registered for %d", message.ID)
	}
}

func (ha *Connection) handleEvent(message Message) {
	var err error
	var timeFired time.Time

	timeFired, err = time.Parse(TimeFormat1, message.Event.TimeFired)
	if err != nil {
		var err2 error
		timeFired, err2 = time.Parse(TimeFormat2, message.Event.TimeFired)
		if err2 != nil {
			fmt.Printf("unable to parse time fired: %+v\n", err)
			fmt.Printf("unable to parse time fired: %+v\n", err2)
			return
		}
	}

	ha.RLock()
	defer ha.RUnlock()

	if handler, ok := ha.handlers[message.ID]; ok {
		go handler.(EventHandler).HandleEvent(ha, message.ID, message.Event.Origin, timeFired, &message.Event)
	} else {
		log.Printf("no handler registered for %d", message.ID)
	}
}
