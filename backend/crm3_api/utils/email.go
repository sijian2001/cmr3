package utils

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// GetEmailConfig メール設定を取得（本番環境では環境変数から読み込む）
func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		Host:     "localhost", // 本番環境ではSMTPサーバーを指定
		Port:     587,
		Username: "your-email@example.com",
		Password: "your-password",
		From:     "noreply@crm3.com",
	}
}

// SendPasswordResetEmail パスワードリセットメールを送信
func SendPasswordResetEmail(to, resetToken string) error {
	config := GetEmailConfig()

	m := gomail.NewMessage()
	m.SetHeader("From", config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "パスワードリセットのご案内")

	// メール本文（HTML）
	resetURL := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", resetToken)
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>パスワードリセットのご案内</h2>
			<p>パスワードリセットのリクエストを受け付けました。</p>
			<p>以下のリンクをクリックして、新しいパスワードを設定してください：</p>
			<p><a href="%s">パスワードをリセットする</a></p>
			<p>このリンクは1時間で期限切れになります。</p>
			<p>もしこのリクエストに覚えがない場合は、このメールを無視してください。</p>
		</body>
		</html>
	`, resetURL)

	m.SetBody("text/html", htmlBody)

	// 開発環境では実際にメールを送信しない（ログのみ出力）
	if config.Host == "localhost" {
		log.Printf("Password reset email would be sent to: %s", to)
		log.Printf("Reset URL: %s", resetURL)
		log.Printf("Email content: %s", htmlBody)
		return nil
	}

	// 本番環境ではSMTP送信
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Password reset email sent to: %s", to)
	return nil
}