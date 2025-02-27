package google

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/api/androidpublisher/v3"
)

type SubscriptionHandler interface {
	OnSubscriptionStarted(ctx context.Context, message *SubscriptionMessage, purchase *androidpublisher.SubscriptionPurchase) error
	OnSubscriptionExpired(ctx context.Context, message *SubscriptionMessage, purchase *androidpublisher.SubscriptionPurchase) error
}

func NewSimpleSubscriptionConsumer(handler SubscriptionHandler) SubscriptionMessageConsumer {
	return &simpleSubscriptionConsumer{
		SubscriptionHandler: handler,
	}
}

type simpleSubscriptionConsumer struct {
	SubscriptionHandler SubscriptionHandler
}

func (s *simpleSubscriptionConsumer) ConsumeMessage(ctx context.Context, message *SubscriptionMessage) bool {
	ds := message.Subscription
	notification := ds.SubscriptionNotification
	if notification == nil {
		return false
	}
	subscriptionInfo, err := FetchSubscription(ctx, ds.PackageName, notification.SubscriptionId, notification.PurchaseToken)
	if err != nil {
		zap.L().Warn("[simpleSubscriptionConsumer] FetchSubscription failed", zap.Error(err))
		return false
	}
	switch notification.NotificationType {
	case SubscriptionNotificationTypePurchased:
		fallthrough
	case SubscriptionNotificationTypeRenewed:
		err = AcknowledgeSubscription(ctx, ds.PackageName, notification.SubscriptionId, notification.PurchaseToken)
		if err != nil {
			zap.L().Error("[simpleSubscriptionConsumer] AcknowledgeSubscription failed", zap.Error(err))
			return false
		}
		err = s.SubscriptionHandler.OnSubscriptionStarted(ctx, message, subscriptionInfo)
		if err != nil {
			zap.L().Warn("[simpleSubscriptionConsumer] OnSubscriptionStarted failed", zap.Error(err))
			return false
		}
		return true
	case SubscriptionNotificationTypeExpired:
		err = s.SubscriptionHandler.OnSubscriptionExpired(ctx, message, subscriptionInfo)
		if err != nil {
			zap.L().Warn("[simpleSubscriptionConsumer] OnSubscriptionExpired failed", zap.Error(err))
			return false
		}
		return true
	case SubscriptionNotificationTypeRecovered:
	case SubscriptionNotificationTypeCanceled:
	case SubscriptionNotificationTypeOnHold:
	case SubscriptionNotificationTypeInGracePeriod:
	case SubscriptionNotificationTypeRestarted:
	case SubscriptionNotificationTypePriceChangeConfirmed:
	case SubscriptionNotificationTypeDeferred:
	case SubscriptionNotificationTypePaused:
	case SubscriptionNotificationTypePauseScheduleChanged:
	case SubscriptionNotificationTypeRevoked:
	}
	return true
}
