package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

type SubscriptionMessage struct {
	Id           string
	PublishTime  time.Time
	Subscription *DeveloperSubscription
}

type MessageConnector interface {
	Attach(ctx context.Context, consumer SubscriptionMessageConsumer) error
}

type SubscriptionMessageConsumer interface {
	ConsumeMessage(ctx context.Context, message *SubscriptionMessage) (ack bool)
}

type pubSubMessageConnector struct {
	ProjectId      string
	SubscriptionId string
}

func NewPubSubMessageConnector(projectId string, subscriptionId string) MessageConnector {
	return &pubSubMessageConnector{
		ProjectId:      projectId,
		SubscriptionId: subscriptionId,
	}
}

func (c *pubSubMessageConnector) Attach(ctx context.Context, consumer SubscriptionMessageConsumer) error {
	client, err := pubsub.NewClient(ctx, c.ProjectId)
	if err != nil {
		return err
	}
	return client.Subscription(c.SubscriptionId).Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		ds := &DeveloperSubscription{}
		e := json.Unmarshal(message.Data, ds)
		if e != nil {
			zap.L().Error("Unmarshal DeveloperSubscription failed", zap.Error(e))
			return
		}
		ack := consumer.ConsumeMessage(ctx, &SubscriptionMessage{
			Id:           message.ID,
			PublishTime:  message.PublishTime,
			Subscription: ds,
		})
		if ack {
			message.Ack()
		} else {
			message.Nack()
		}
	})
}
