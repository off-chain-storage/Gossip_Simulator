package iface

type ReadOnlyRedisDB interface {
	GetDataFromRedis(key string) (string, error)
}

type AccessRedisDB interface {
	ReadOnlyRedisDB

	SetDataToRedis(key, value string) error
	// DeleteDataFromRedis(key string) error
}

type RedisDB interface {
	AccessRedisDB

	SetRedisConn()
}
