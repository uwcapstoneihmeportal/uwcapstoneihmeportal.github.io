package handlers
import (
"log"

"github.com/gorilla/websocket"
"sync"
	"encoding/json"
)

//Notifier is an object that handles WebSocket notifications.
//It starts a new
type Notifier struct {
	//eventQ is a go channel that
	//into which one goroutine can
	//write byte slices, and out of which
	//another goroutine can read those byte slices
	eventQ chan []byte

	//TODO: add other fields to this struct
	//that you might need. For example, you'll
	//need to track all of the current WebSocket
	//connections. Remember that slice will be used
	//by multiple goroutines, so you'll need to
	//protect it for concurrent use!
	targetUsers []int64
	conns map[int64]map[*websocket.Conn] struct{}
	//in the hwk use conns map[int64] map[*websocket.Conn]struct{}
	mx sync.Mutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := &Notifier{
		eventQ: make(chan []byte, 1024), //buffered channel that can hold 1024 slices at a time
		conns: make(map[int64]map[*websocket.Conn]struct{}),
		targetUsers: []int64{},
	}
	go n.start()
	//TODO: call the .start() method on its own goroutine

	return n
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(userID int64, client *websocket.Conn) {
	log.Println("adding new WebSockets client")
	//TODO: add the client to the slice you are using
	//to track all current WebSocket connections.
	//Since this can be called from multiple
	//goroutines, make sure you protect that slice
	//while you add a new connection to it!
	n.mx.Lock()
	_, exists := n.conns[userID]
	if !exists {
		n.conns[userID] = map[*websocket.Conn]struct{}{}
	}
	n.conns[userID][client] = struct{}{}
	n.mx.Unlock()
	//also process incoming control messages from
	//the client, as described in this section of the docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
	for {
		if _,_,err := client.NextReader(); err != nil {
			client.Close()
			break
		}
	}
}

//Notify broadcasts the event to all WebSocket clients
func (n *Notifier) Notify(event []byte) {
	log.Printf("adding event to the queue")
	//TODO: add `event` to the `n.eventQ`
	//see https://tour.golang.org/concurrency/2
	//and https://gobyexample.com/channels
	n.eventQ <- event[:]
}

type MQUserIDs struct{
	UserIDs []int64 `json:"userIDs"`
}

//start starts the notification loop
func (n *Notifier) start() {
	log.Println("starting notifier loop")

	//TODO: start a never-ending loop that reads
	//new events out of the `n.eventQ` and broadcasts
	//them to all WebSocket connections.
	//To write the byte-slice to the WebSocket, use
	//the .WriteMessage() method.
	//https://godoc.org/github.com/gorilla/websocket#Conn.WriteMessage
	//Or, for better performance, prepare the message once
	//and use the .WritePreparedMessage() method.
	//https://godoc.org/github.com/gorilla/websocket#PreparedMessage
	for ev := range n.eventQ {
		prepMessage, err := websocket.NewPreparedMessage(websocket.TextMessage, ev)
		if err != nil {
			log.Println(err)
		}
		var userIds MQUserIDs
		json.Unmarshal(ev, userIds)
		log.Println(userIds.UserIDs)
		n.mx.Lock()
		if len(userIds.UserIDs) == 0 {
			log.Println(n.conns)
			for userID, kmp := range n.conns {
				for k := range kmp {
					if err = k.WritePreparedMessage(prepMessage); err != nil {
						delete(n.conns[userID], k)
					}
				}
			}
		} else {
			for userID := range userIds.UserIDs {
				for k := range n.conns[int64(userID)] {
					if err = k.WritePreparedMessage(prepMessage); err != nil {
						delete(n.conns[int64(userID)], k)
					}
				}
			}
		}
		n.mx.Unlock()
	}
	//Remember that you need to lock the slice of connections
	//while you iterate it, as other goroutines might
	//be trying to add new clients to it while you iterate!

	//If you get an error while trying to write the
	//message to one of the WebSocket connections,
	//that means the client has disconnected, so
	//remove that connection from your list.
}
