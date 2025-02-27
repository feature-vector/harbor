package google

import "context"

type chainedConsumer struct {
	Consumers []SubscriptionMessageConsumer
}

func NewChainedConsumer(consumers ...SubscriptionMessageConsumer) SubscriptionMessageConsumer {
	return &chainedConsumer{Consumers: consumers}
}

func (c *chainedConsumer) ConsumeMessage(ctx context.Context, message *SubscriptionMessage) bool {
	for _, consumer := range c.Consumers {
		ack := consumer.ConsumeMessage(ctx, message)
		if ack {
			return true
		}
	}
	return false
}
