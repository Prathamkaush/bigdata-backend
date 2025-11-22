package metrics

import "sync/atomic"

// Request counters
var ApiRequestsTotal uint64
var ApiRequests2xx uint64
var ApiRequests4xx uint64
var ApiRequests5xx uint64
var RateLimitBlocked uint64

// Cache counters
var CacheHits uint64
var CacheMisses uint64

//rate limiting
var RateLimitMax uint64 = 10
var RateLimitUsed uint64 = 0

// Helper functions (optional)
func IncHits() {
	atomic.AddUint64(&CacheHits, 1)
}

func IncMiss() {
	atomic.AddUint64(&CacheMisses, 1)
}
func RateLimitRemaining() uint64 {
	if RateLimitMax <= RateLimitUsed {
		return 0
	}
	return RateLimitMax - RateLimitUsed
}
