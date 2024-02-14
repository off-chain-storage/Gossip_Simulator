package db

import "flag-example/curie-node/db/iface"

// DB에 접근하여 오직 Data Read만을 지원하는 인터페이스
type ReadOnlyRedisDB = iface.ReadOnlyRedisDB

// DB에 접근하여 Data Read/Write를 지원하는 인터페이스
type AccessRedisDB = iface.AccessRedisDB

// 전체 데이터베이스 인터페이스
type RedisDB = iface.RedisDB
