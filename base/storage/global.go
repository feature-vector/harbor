package storage

var (
	globalProvider StorageProvider
)

func SetGlobalProvider(provider StorageProvider) {
	globalProvider = provider
}

func Global() StorageProvider {
	return globalProvider
}
