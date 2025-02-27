package google

type OneTimeProductNotificationType int
type SubscriptionNotificationType int

// OneTimeProductNotificationType
// (1) ONE_TIME_PRODUCT_PURCHASED - 用户成功购买了一次性商品。
// (2) ONE_TIME_PRODUCT_CANCELED - 用户已取消待处理的一次性商品购买交易。
const (
	OneTimeProductNotificationTypePurchased OneTimeProductNotificationType = 1
	OneTimeProductNotificationTypeCanceled  OneTimeProductNotificationType = 2
)

// SubscriptionNotificationType
// (1) SUBSCRIPTION_RECOVERED - 从帐号保留状态恢复了订阅。
// (2) SUBSCRIPTION_RENEWED - 续订了处于活动状态的订阅。
// (3) SUBSCRIPTION_CANCELED - 自愿或非自愿地取消了订阅。如果是自愿取消，在用户取消时发送。
// (4) SUBSCRIPTION_PURCHASED - 购买了新的订阅。
// (5) SUBSCRIPTION_ON_HOLD - 订阅已进入帐号保留状态（如果已启用）。
// (6) SUBSCRIPTION_IN_GRACE_PERIOD - 订阅已进入宽限期（如果已启用）。
// (7) SUBSCRIPTION_RESTARTED - 用户已通过 Play > 帐号 > 订阅恢复了订阅。订阅已取消，但在用户恢复时尚未到期。如需了解详情，请参阅 [恢复](/google/play/billing/subscriptions#restore)。
// (8) SUBSCRIPTION_PRICE_CHANGE_CONFIRMED - 用户已成功确认订阅价格变动。
// (9) SUBSCRIPTION_DEFERRED - 订阅的续订时间点已延期。
// (10) SUBSCRIPTION_PAUSED - 订阅已暂停。
// (11) SUBSCRIPTION_PAUSE_SCHEDULE_CHANGED - 订阅暂停计划已更改。
// (12) SUBSCRIPTION_REVOKED - 用户在到期时间之前已撤消订阅。
// (13) SUBSCRIPTION_EXPIRED - 订阅已到期。
const (
	SubscriptionNotificationTypeRecovered            SubscriptionNotificationType = 1
	SubscriptionNotificationTypeRenewed              SubscriptionNotificationType = 2
	SubscriptionNotificationTypeCanceled             SubscriptionNotificationType = 3
	SubscriptionNotificationTypePurchased            SubscriptionNotificationType = 4
	SubscriptionNotificationTypeOnHold               SubscriptionNotificationType = 5
	SubscriptionNotificationTypeInGracePeriod        SubscriptionNotificationType = 6
	SubscriptionNotificationTypeRestarted            SubscriptionNotificationType = 7
	SubscriptionNotificationTypePriceChangeConfirmed SubscriptionNotificationType = 8
	SubscriptionNotificationTypeDeferred             SubscriptionNotificationType = 9
	SubscriptionNotificationTypePaused               SubscriptionNotificationType = 10
	SubscriptionNotificationTypePauseScheduleChanged SubscriptionNotificationType = 11
	SubscriptionNotificationTypeRevoked              SubscriptionNotificationType = 12
	SubscriptionNotificationTypeExpired              SubscriptionNotificationType = 13
)

type DeveloperSubscription struct {
	Version                    string                      `json:"version"`
	PackageName                string                      `json:"packageName"`
	EventTimeMillis            string                      `json:"eventTimeMillis"`
	OneTimeProductNotification *OneTimeProductNotification `json:"oneTimeProductNotification"`
	SubscriptionNotification   *SubscriptionNotification   `json:"subscriptionNotification"`
	TestNotification           *TestNotification           `json:"testNotification"`
	VoidedPurchaseNotification *VoidedPurchaseNotification `json:"voidedPurchaseNotification"`
}

type VoidedPurchaseNotification struct {
	PurchaseToken string `json:"purchaseToken"`
	OrderId       string `json:"orderId"`
	ProductType   int    `json:"productType"`
	RefundType    int    `json:"refundType"`
}

type OneTimeProductNotification struct {
	Version          string                         `json:"version"`
	NotificationType OneTimeProductNotificationType `json:"notificationType"`
	PurchaseToken    string                         `json:"purchaseToken"`
	Sku              string                         `json:"sku"`
}

type SubscriptionNotification struct {
	Version          string                       `json:"version"`
	NotificationType SubscriptionNotificationType `json:"notificationType"`
	PurchaseToken    string                       `json:"purchaseToken"`
	SubscriptionId   string                       `json:"subscriptionId"`
}

type TestNotification struct {
	Version string `json:"version"`
}
