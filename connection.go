package homeassistant

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Connection .
type Connection struct {
	sync.RWMutex
	conn       *websocket.Conn
	handlers   map[uint64]interface{}
	lastID     uint64
	Done       chan struct{}
	Authorized chan bool
}

// NewConnection .
func NewConnection(hostname string, port int, accessToken string, secure bool) (*Connection, error) {
	ha := new(Connection)

	var u url.URL
	if secure {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Host = hostname
	if port == 0 {
		u.Host = hostname
	} else {
		u.Host = fmt.Sprintf("%s:%d", hostname, port)
	}
	u.Path = "/api/websocket"

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	ha.conn = c
	ha.handlers = make(map[uint64]interface{})
	ha.Done = make(chan struct{})
	ha.Authorized = make(chan bool)

	go func() {
		defer ha.conn.Close()
		defer close(ha.Done)

		for {
			var message Message
			err := c.ReadJSON(&message)
			if err != nil {
				if ce, ok := err.(*websocket.CloseError); ok {
					switch ce.Code {
					case websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
						websocket.CloseNoStatusReceived:
						//log.Printf("normal closure")
						return
					}
				}
				log.Printf("read: %T %v\n", err, err)
				return
			}
			log.Printf("recv: %v", message)

			switch message.Type {
			case "auth_required":
				var message = struct {
					Type        string `json:"type"`
					AccessToken string `json:"access_token"`
				}{
					Type:        "auth",
					AccessToken: accessToken,
				}

				err = c.WriteJSON(&message)
				if err != nil {
					log.Fatal("write:", err)
				}

			case "auth_ok":
				ha.Authorized <- true
				fmt.Printf("auth ok\n")

			case "auth_invalid":
				ha.Authorized <- false
				log.Printf("auth failure: %s", message.Message)

			case "result":
				ha.handleResult(message)

			case "event":
				ha.handleEvent(message)

			case "pong":
				ha.handlePong(message)
				log.Printf("pong\n")

			default:
				log.Printf("some other type\n")
			}
		}
	}()

	return ha, nil
}

func (ha *Connection) sendMessage(handler interface{}, message Command) (uint64, error) {

	ha.Lock()
	defer ha.Unlock()

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
	log.Printf("closing\n")
	err := ha.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			""),
	)
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

	timeFired, err := time.Parse(TimeFormat, message.Event.TimeFired)
	if err != nil {
		fmt.Printf("unable to parse time fired: %s", err.Error())
		return
	}

	fmt.Println(message.Event.String())

	ha.RLock()
	defer ha.RUnlock()
	if handler, ok := ha.handlers[message.ID]; ok {
		go handler.(EventHandler).HandleEvent(ha, message.ID, message.Event.Origin, timeFired, &message.Event)
	} else {
		log.Printf("no handler registered for %d", message.ID)
	}
}
