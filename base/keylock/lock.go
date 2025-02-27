package keylock

type KeyLocker interface {
	Lock(key string)
	Unlock(key string)
	TryLock(key string) bool
}
