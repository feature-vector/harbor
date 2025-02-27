package google

import (
	"context"
	"google.golang.org/api/androidpublisher/v3"
)

func FetchProduct(
	ctx context.Context,
	packageName string,
	productID string,
	purchaseToken string,
) (*androidpublisher.ProductPurchase, error) {
	service, err := androidpublisher.NewService(ctx)
	if err != nil {
		return nil, err
	}

	ps := androidpublisher.NewPurchasesProductsService(service)
	result, err := ps.Get(packageName, productID, purchaseToken).Context(ctx).Do()

	return result, err
}

func AcknowledgeProduct(
	ctx context.Context,
	packageName string,
	sku string,
	purchaseToken string,
) error {
	service, err := androidpublisher.NewService(ctx)
	if err != nil {
		return err
	}

	ps2 := androidpublisher.NewPurchasesProductsService(service)
	return ps2.Acknowledge(packageName, sku, purchaseToken, &androidpublisher.ProductPurchasesAcknowledgeRequest{}).Context(ctx).Do()
}

func ConsumeProduct(
	ctx context.Context,
	packageName string,
	sku string,
	purchaseToken string,
) error {
	service, err := androidpublisher.NewService(ctx)
	if err != nil {
		return err
	}

	ps2 := androidpublisher.NewPurchasesProductsService(service)
	return ps2.Consume(packageName, sku, purchaseToken).Context(ctx).Do()
}

func FetchSubscription(
	ctx context.Context,
	packageName string,
	subscriptionId string,
	purchaseToken string,
) (*androidpublisher.SubscriptionPurchase, error) {
	service, err := androidpublisher.NewService(ctx)
	if err != nil {
		return nil, err
	}

	ps2 := androidpublisher.NewPurchasesSubscriptionsService(service)
	return ps2.Get(packageName, subscriptionId, purchaseToken).Context(ctx).Do()
}

func AcknowledgeSubscription(
	ctx context.Context,
	packageName string,
	subscriptionId string,
	purchaseToken string,
) error {
	service, err := androidpublisher.NewService(ctx)
	if err != nil {
		return err
	}

	ps2 := androidpublisher.NewPurchasesSubscriptionsService(service)
	return ps2.Acknowledge(packageName, subscriptionId, purchaseToken, &androidpublisher.SubscriptionPurchasesAcknowledgeRequest{}).Context(ctx).Do()
}
