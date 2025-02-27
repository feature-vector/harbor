package google

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/api/androidpublisher/v3"
)

type ProductHandler interface {
	OnProductPurchased(ctx context.Context, message *SubscriptionMessage, purchase *androidpublisher.ProductPurchase) error
}

func NewSimpleProductConsumer(handler ProductHandler) SubscriptionMessageConsumer {
	return &simpleProductConsumer{
		ProductHandler: handler,
	}
}

type simpleProductConsumer struct {
	ProductHandler ProductHandler
}

func (s *simpleProductConsumer) ConsumeMessage(ctx context.Context, message *SubscriptionMessage) bool {
	ds := message.Subscription
	notification := ds.OneTimeProductNotification
	if notification == nil {
		return false
	}
	productInfo, err := FetchProduct(ctx, ds.PackageName, notification.Sku, notification.PurchaseToken)
	if err != nil {
		zap.L().Warn("[simpleProductConsumer] FetchProduct failed", zap.Error(err))
		return false
	}
	if productInfo.PurchaseState != 0 {
		return true
	}
	err = AcknowledgeProduct(ctx, ds.PackageName, notification.Sku, notification.PurchaseToken)
	if err != nil {
		zap.L().Warn("[simpleProductConsumer] AcknowledgeProduct failed", zap.Error(err))
		return false
	}
	err = ConsumeProduct(ctx, ds.PackageName, notification.Sku, notification.PurchaseToken)
	if err != nil {
		zap.L().Warn("[simpleProductConsumer] ConsumeProduct failed", zap.Error(err))
		return false
	}
	err = s.ProductHandler.OnProductPurchased(ctx, message, productInfo)
	return err == nil
}
