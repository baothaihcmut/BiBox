package models

import "time"

type NotificationSession struct {
	UserId   string
	ExpireAt time.Time
}
