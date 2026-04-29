package db

import (
	"database/sql"
	"hash/crc32"
)

type ShardManager struct {
	Shards map[int]*sql.DB
}

func (s *ShardManager) getShardID(key string) int {
	hash := crc32.ChecksumIEEE([]byte(key))
	return int(hash % 2)
}

func (s *ShardManager) GetDB(key string) *sql.DB {
	shardID := s.getShardID(key)
	return s.Shards[shardID]
}
