package kafka

import (
	"context"
	"fmt"
	"marketing/internal/transport/kafka/handler"
	"marketing/internal/transport/kafka/message"
)

type Handler interface {
	Handle(ctx context.Context, msg message.Message) error
}

type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

func (r *Router) Register(topic string, handler Handler) {
	r.routes[topic] = handler
}

func (r *Router) Route(ctx context.Context, msg message.Message) error {
	routeHandler, exists := r.routes[msg.Topic]
	if !exists {
		return fmt.Errorf("no handler registered for topic: %s", msg.Topic)
	}

	return routeHandler.Handle(ctx, msg)
}

func (r *Router) GetTopics() []string {
	topics := make([]string, 0, len(r.routes))

	for topic := range r.routes {
		topics = append(topics, topic)
	}

	return topics
}

func getRouter(logger Logger) *Router {
	r := NewRouter()

	r.Register("user.add", handler.NewUserCreatedHandler(logger))

	return r
}
