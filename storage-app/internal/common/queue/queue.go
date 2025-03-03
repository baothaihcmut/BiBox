package queue

type QueueService interface {
	PublishMessage(topic string, value interface{}, headers map[string]string) (int32, int64, error)
}
