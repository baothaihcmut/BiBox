package response

import (
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
)

type NotificationOutput struct {
	Id         string               `json:"id"`
	UserId     string               `json:"user_id"`
	Type       enums.NoficationType `json:"type"`
	Title      string               `json:"title"`
	Message    string               `json:"message"`
	ActionUrl  string               `json:"action_url"`
	Seen       bool                 `json:"seen"`
	FromUserId string               `json:"from_user_id"`
	CreatedAt  time.Time            `json:"created_at"`
}
