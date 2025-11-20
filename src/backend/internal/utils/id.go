package utils

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ulidEntropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	ulidMutex   sync.Mutex
)

// GenerateULID 生成全局唯一的ULID字符串
func GenerateULID() string {
	ulidMutex.Lock()
	defer ulidMutex.Unlock()
	id := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	return id.String()
}
