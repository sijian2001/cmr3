package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcryptでも空文字はハッシュ可能
		},
		{
			name:     "Long password (72 bytes limit)",
			password: "this_is_a_very_long_password_with_more_than_72_characters_to_test_the_bcrypt_limit_properly",
			wantErr:  true, // bcryptは72バイト以上のパスワードでエラーを返す
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// ハッシュが生成されていることを確認
				if len(hash) == 0 {
					t.Errorf("HashPassword() returned empty hash")
				}

				// 元のパスワードと異なることを確認
				if hash == tt.password {
					t.Errorf("HashPassword() returned unhashed password")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "test_password_123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Incorrect password",
			password: "wrong_password",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalid_hash",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}