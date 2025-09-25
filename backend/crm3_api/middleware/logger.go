package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestResponseLogger struct {
	next echo.HandlerFunc
}

func RequestResponseLogging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()

			// リクエストボディを読み取り（認証APIの場合は機密情報をマスク）
			var reqBody []byte
			if req.Body != nil {
				reqBody, _ = io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			// レスポンスを記録するためのカスタムレスポンスライター
			resBody := new(bytes.Buffer)
			writer := &responseBodyWriter{ResponseWriter: res.Writer, body: resBody}
			res.Writer = writer

			// リクエストをログ出力（パスワードなどの機密情報をマスク）
			maskedReqBody := maskSensitiveData(reqBody, req.URL.Path)
			log.Printf("REQUEST: %s %s | Body: %s | IP: %s | UserAgent: %s",
				req.Method,
				req.URL.Path,
				string(maskedReqBody),
				c.RealIP(),
				req.Header.Get("User-Agent"),
			)

			// 次のハンドラーを実行
			err := next(c)

			// 処理時間計算
			latency := time.Since(start)

			// レスポンスをログ出力
			log.Printf("RESPONSE: %s %s | Status: %d | Latency: %v | Size: %d | Body: %s",
				req.Method,
				req.URL.Path,
				res.Status,
				latency,
				res.Size,
				resBody.String(),
			)

			return err
		}
	}
}

type responseBodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// maskSensitiveData は機密情報をマスクする
func maskSensitiveData(body []byte, path string) []byte {
	if len(body) == 0 {
		return body
	}

	// 認証関連のエンドポイントの場合、パスワードをマスク
	if strings.Contains(path, "/auth/") {
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return body
		}

		// パスワード関連フィールドをマスク
		if _, exists := data["password"]; exists {
			data["password"] = "***MASKED***"
		}
		if _, exists := data["confirm_password"]; exists {
			data["confirm_password"] = "***MASKED***"
		}
		if _, exists := data["old_password"]; exists {
			data["old_password"] = "***MASKED***"
		}

		maskedBody, err := json.Marshal(data)
		if err != nil {
			return body
		}
		return maskedBody
	}

	return body
}