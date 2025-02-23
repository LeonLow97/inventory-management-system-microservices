package cache

import (
	"context"
	"fmt"
	"log"
	"time"
)

func lockKey(userID int, bucketName string) string {
	return fmt.Sprintf("lock:%d:%s", userID, bucketName)
}

// AcquireLock prevents race conditions in updating tokens in bucket
// Uses Redis SetNX command to acquire a lock for a user and a specific bucket
func (c Cache) AcquireLock(ctx context.Context, userID int, bucketName string) bool {
	lockKey := lockKey(userID, bucketName)
	lockExpiration := time.Second * time.Duration(c.cfg.RateLimiting.BucketLockExpiration)
	success, err := c.RedisClient.SetNX(ctx, lockKey, 1, lockExpiration).Result()
	if err != nil {
		log.Printf("failed to acquire distributed lock for userid %d\n", userID)
		return false
	}
	return success
}

// ReleaseLock releases the distributed lock after the operation
// allowing other replicas to access the rate limiting logic
func (c Cache) ReleaseLock(ctx context.Context, userID int, bucketName string) {
	lockKey := lockKey(userID, bucketName)
	c.RedisClient.Del(ctx, lockKey)
}
