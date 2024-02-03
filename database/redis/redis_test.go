package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisConnectionWithMock(t *testing.T) {
	// Create a miniredis server with auto-cleanup
	mr := miniredis.RunT(t)

	r := New(&Config{Addr: mr.Addr()})

	_, err := r.Set("slug", "https://redis.io/docs/connect/clients/go/", 0).Result()

	assert.NoError(t, err, "Failed to set cache key")

	result, err := r.Get("slug").Result()
	assert.NoError(t, err, "Failed to get value from redis")

	assert.Equal(t, "https://redis.io/docs/connect/clients/go/", result)
}
