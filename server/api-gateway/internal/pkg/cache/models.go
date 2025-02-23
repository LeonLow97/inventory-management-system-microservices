package cache

import "fmt"

func (c Cache) UserTokenBucketKey(userID int, bucketName string) string {
	return fmt.Sprintf("user:%d:%s", userID, bucketName)
}
