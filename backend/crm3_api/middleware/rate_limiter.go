package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimiterConfig レートリミッター設定
type RateLimiterConfig struct {
	Requests int           // 許可されるリクエスト数
	Window   time.Duration // 時間ウィンドウ
	KeyFunc  func(c echo.Context) string // クライアント識別用のキー生成関数
}

// ClientInfo クライアントのリクエスト情報
type ClientInfo struct {
	Count     int
	ResetTime time.Time
}

// RateLimiter レートリミッター構造体
type RateLimiter struct {
	clients map[string]*ClientInfo
	config  RateLimiterConfig
	mutex   sync.Mutex
}

// NewRateLimiter 新しいレートリミッターを作成
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.KeyFunc == nil {
		config.KeyFunc = func(c echo.Context) string {
			return c.RealIP()
		}
	}

	rl := &RateLimiter{
		clients: make(map[string]*ClientInfo),
		config:  config,
	}

	// 定期的にクライアント情報をクリーンアップ
	go rl.cleanup()

	return rl
}

// RateLimit レートリミット用のミドルウェア関数
func (rl *RateLimiter) RateLimit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := rl.config.KeyFunc(c)

			rl.mutex.Lock()
			defer rl.mutex.Unlock()

			now := time.Now()
			client, exists := rl.clients[key]

			if !exists || now.After(client.ResetTime) {
				// 新しいクライアントまたは時間ウィンドウがリセット
				rl.clients[key] = &ClientInfo{
					Count:     1,
					ResetTime: now.Add(rl.config.Window),
				}
				return next(c)
			}

			if client.Count >= rl.config.Requests {
				// レート制限に達している
				c.Response().Header().Set("X-RateLimit-Limit", string(rune(rl.config.Requests)))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")
				c.Response().Header().Set("X-RateLimit-Reset", string(rune(client.ResetTime.Unix())))

				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error":   "rate_limit_exceeded",
					"message": "リクエスト制限に達しました。しばらく時間をおいてから再度お試しください。",
					"retry_after": int(client.ResetTime.Sub(now).Seconds()),
				})
			}

			// リクエストをカウント
			client.Count++
			c.Response().Header().Set("X-RateLimit-Limit", string(rune(rl.config.Requests)))
			c.Response().Header().Set("X-RateLimit-Remaining", string(rune(rl.config.Requests-client.Count)))
			c.Response().Header().Set("X-RateLimit-Reset", string(rune(client.ResetTime.Unix())))

			return next(c)
		}
	}
}

// cleanup 期限切れのクライアント情報を定期的にクリーンアップ
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute * 5) // 5分毎にクリーンアップ
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mutex.Lock()
			now := time.Now()
			for key, client := range rl.clients {
				if now.After(client.ResetTime) {
					delete(rl.clients, key)
				}
			}
			rl.mutex.Unlock()
		}
	}
}

// 事前定義されたレートリミッター設定

// AuthRateLimiter 認証用レートリミッター（厳しい制限）
func AuthRateLimiter() echo.MiddlewareFunc {
	limiter := NewRateLimiter(RateLimiterConfig{
		Requests: 5,                // 5リクエスト
		Window:   time.Minute * 15, // 15分間
	})
	return limiter.RateLimit()
}

// GeneralRateLimiter 一般API用レートリミッター
func GeneralRateLimiter() echo.MiddlewareFunc {
	limiter := NewRateLimiter(RateLimiterConfig{
		Requests: 100,            // 100リクエスト
		Window:   time.Minute,    // 1分間
	})
	return limiter.RateLimit()
}

// PasswordResetRateLimiter パスワードリセット用レートリミッター（非常に厳しい制限）
func PasswordResetRateLimiter() echo.MiddlewareFunc {
	limiter := NewRateLimiter(RateLimiterConfig{
		Requests: 3,                // 3リクエスト
		Window:   time.Hour,        // 1時間
	})
	return limiter.RateLimit()
}