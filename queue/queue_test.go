package queue

import "testing"

func TestRedisOptDefaultAddr(t *testing.T) {
	opt := redisOpt(RedisConfig{})
	if opt.Addr != "127.0.0.1:6379" {
		t.Fatalf("addr = %q", opt.Addr)
	}
}
