package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/services"
	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	Init(r *gin.RouterGroup)
}

type NotifactionControllerImpl struct {
	notificationSSEManager services.NotificationSSEManagerService
}

func (n *NotifactionControllerImpl) Init(r *gin.RouterGroup) {
	sse := r.Group("/notifications/sse")
	sse.GET("", n.handleSSENotification)
}

func (n *NotifactionControllerImpl) handleSSENotification(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	flusher, _ := c.Writer.(http.Flusher)
	sessionId := c.Query("session_id")
	msgCh, userId, err := n.notificationSSEManager.Connect(c.Request.Context(), sessionId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	for msg := range msgCh {
		select {
		case <-c.Request.Context().Done():
			n.notificationSSEManager.Disconnect(c.Request.Context(), userId, sessionId)
			return
		default:
			jsonData, err := json.Marshal(msg)
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}
			c.Writer.Write(jsonData)
			flusher.Flush()
		}
	}

}
