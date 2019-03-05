package handlers

import (
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

type SessionState struct {
	Time time.Time `json:"time"`
	AuthUser *users.User `json:"authUser"`
}
