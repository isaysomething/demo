package users

import "time"

type LoginEvent struct {
	User     *User
	Duration time.Duration
}

type LogoutEvent struct {
	User *User
}
