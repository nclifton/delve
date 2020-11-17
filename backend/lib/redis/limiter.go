package redis

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-redis/redis/v7"
)

const redisLimitScript = `
local tokens_key = KEYS[1]
local timestamp_key = KEYS[2]
local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local fill_time = capacity/rate
local ttl = math.floor(fill_time*2)
local last_tokens = tonumber(redis.call("get", tokens_key))
if last_tokens == nil then
    last_tokens = capacity
end
local last_refreshed = tonumber(redis.call("get", timestamp_key))
if last_refreshed == nil then
    last_refreshed = 0
end
local delta = math.max(0, now-last_refreshed)
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end
redis.call("setex", tokens_key, ttl, new_tokens)
redis.call("setex", timestamp_key, ttl, now)
return { allowed, new_tokens }
`

// Inf is the infinite rate limit; it allows all events (even if burst is zero)
const Inf = float64(math.MaxFloat64)

// Limiter controls how frequently events are allowed to happen
type Limiter struct {
	client     *redis.Client
	scriptHash string
}

// New returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens
func NewLimiter(redisURL string) (*Limiter, error) {
	conn, err := Connect(redisURL)

	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(redisLimitScript)))
	exists, err := conn.Client.ScriptExists(hash).Result()
	if err != nil {
		return nil, err
	}

	// load script when missing
	if !exists[0] {
		_, err := conn.Client.ScriptLoad(redisLimitScript).Result()
		if err != nil {
			return nil, err
		}
	}

	return &Limiter{
		client:     conn.Client,
		scriptHash: hash,
	}, nil
}

// Allow is shorthand for AllowN(time.Now(), 1)
func (l *Limiter) Allow(key string, rate float64, burst int) bool {
	return l.reserveN(time.Now(), 1, key, rate, burst).ok
}

// Reservation contains tokens for use if available
// in future we may allow returning tokens to the bucket
type Reservation struct {
	ok     bool
	tokens int
}

func (l *Limiter) reserveN(now time.Time, n int, key string, rate float64, burst int) Reservation {
	results, err := l.client.EvalSha(
		l.scriptHash,
		[]string{key + ".tokens", key + ".time"},
		rate,
		burst,
		now.Unix(),
		n,
	).Result()
	if err != nil {
		log.Println("failed to call rate limit:", err)
		return Reservation{
			ok:     true,
			tokens: n,
		}
	}

	rs, ok := results.([]interface{})
	if !ok {
		return Reservation{
			ok:     false,
			tokens: n,
		}
	}

	newTokens, _ := rs[1].(int64)
	return Reservation{
		ok:     rs[0] == int64(1),
		tokens: int(newTokens),
	}
}
