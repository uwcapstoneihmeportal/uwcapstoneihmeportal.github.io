package handlers

import (
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/models/users"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/indexes"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/sessions"
	"github.com/gorilla/websocket"
	"net/http"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

type SessionContext struct {
	signInKey string
	sessionStore sessions.Store
	userStore users.MySQLStore
	lookUpTrie *indexes.Trie
	notifier *Notifier
	upgrader websocket.Upgrader
}

func NewContext(signInKey string, store sessions.Store, userStore users.MySQLStore, lookupTrie *indexes.Trie, notifier *Notifier) *SessionContext {
	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &SessionContext{signInKey, store, userStore, lookupTrie, notifier, upgrader}
}
