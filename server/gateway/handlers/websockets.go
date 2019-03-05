package handlers

import (
	"net/http"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/sessions"
)

type RabbitMQMessage struct {
	Body string
	UserIDs []int64
}

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.
func (ctx *SessionContext) WebsocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	var sessState *SessionState
	_, err := sessions.GetState(r, ctx.signInKey, ctx.sessionStore, &sessState)
	if err != nil {
		http.Error(w, fmt.Sprintf("error in getting state: %v", err), http.StatusUnauthorized)
		return
	}
	conn, err := ctx.upgrader.Upgrade(w,r,nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error in upgrading the client to websocket connection: %v", err), http.StatusInternalServerError)
	}
	ctx.notifier.AddClient(sessState.AuthUser.UserID ,conn)
}

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

func RabbitMQSubscriber(notifier *Notifier, conn *amqp.Connection) {
	mqChan, err := conn.Channel()
	reportError(err, "Error in creating channel to connecto to RabbitMQ")

	q, err := mqChan.QueueDeclare("eventQ",
		false,
			false,
				false,
					false,
						nil)
	reportError(err, "Error in declaring a queue for mq Channel")
	messages, err := mqChan.Consume(q.Name,
		"",
			false,
				false,
					false,
						false,
							nil)
	go processMessages(messages, notifier)
}

func processMessages(msgs <- chan amqp.Delivery, notifier *Notifier) {
	for msg := range msgs {
		log.Printf("received message: %v", string(msg.Body))
		notifier.Notify(msg.Body)
		msg.Ack(false)
	}
}

func reportError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}