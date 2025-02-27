# subscription问题


## subscription ack的问题
* 收到google cloud的推送需要ACK，这个是google cloud本身的ACK，不ACK，google cloud还会给你推事件。
* 如果是新支付的subscription，需要在googleapi上给这个subscription ACK（续费的不需要ACK），表示收到了这个subscription。这个是googleapi的ACK，和google cloud的ACK无关。
* 如果一定时间没有调用googleapi进行ACK，这个subscription会被系统取消。

## orderId的问题
* subscription的orderId不是固定的，每次续费时的orderId后面会有..0/..1/..2/..3这种后缀。
* 前面的orderId是否唯一不确定，有人建议用purchaseToken唯一标识一个订单


## LinkedPurchaseToken和ObfuscatedExternalAccountId
* 如果用ObfuscatedExternalAccountId标识支付的用户id，会有一些问题要处理
* 如果用户升级or降级，那么会重新支付，ObfuscatedExternalAccountId会是空的，LinkedPurchaseToken会标识到之前subscription的PurchaseToken
* 在google商店某些未知的购买行为，也有LinkedPurchaseToken和ObfuscatedExternalAccountId都为空的情况，只能通过恢复购买获取权益





