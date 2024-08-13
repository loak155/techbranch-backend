package mock

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/loak155/techbranch-backend/pkg/redis"
)

func NewRedisMock(t *testing.T, db int, expiration time.Duration) *redis.RedisManager {
	t.Helper()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("unexpected error while createing test redis server '%#v'", err)
	}

	return redis.NewRedisManager(s.Addr(), db, expiration)
}
