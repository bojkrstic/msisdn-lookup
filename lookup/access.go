package lookup

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	apiKeyHeader       = "X-API-Key"
	defaultAPIKeyValue = "local-dev-key"
	rateLimitPerMinute = 60
)

type rateBucket struct {
	Count  int
	Window time.Time
}

var (
	apiKeyOnce sync.Once
	apiKey     string

	rateMutex sync.Mutex
	rateStore = make(map[string]*rateBucket)
)

func enforceAccess(w http.ResponseWriter, r *http.Request) bool {
	requiredKey := DefaultAPIKey()
	provided := r.Header.Get(apiKeyHeader)
	if provided == "" || provided != requiredKey {
		writeJSONError(w, http.StatusUnauthorized, map[string]string{"error": "invalid API key"})
		return false
	}

	if !allowKey(provided) {
		writeJSONError(w, http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
		return false
	}

	return true
}

func DefaultAPIKey() string {
	apiKeyOnce.Do(func() {
		apiKey = os.Getenv("LOOKUP_API_KEY")
		if apiKey == "" {
			apiKey = defaultAPIKeyValue
		}
	})
	return apiKey
}

func allowKey(key string) bool {
	rateMutex.Lock()
	defer rateMutex.Unlock()

	bucket, ok := rateStore[key]
	if !ok {
		bucket = &rateBucket{Window: time.Now()}
		rateStore[key] = bucket
	}

	now := time.Now()
	if now.Sub(bucket.Window) >= time.Minute {
		bucket.Window = now
		bucket.Count = 0
	}

	if bucket.Count >= rateLimitPerMinute {
		return false
	}

	bucket.Count++
	return true
}

func writeJSONError(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
