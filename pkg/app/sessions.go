package app

import "time"

const sessionKey = "session"

var sessions = make(map[string]sessionCookie)

type sessionCookie struct {
	Name    string
	Value   string
	Expires time.Time
	MaxAge  int
}
