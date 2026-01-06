package main

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	userEmailKey contextKey = "userEmail"
)

// AuthMiddleware は IAP ヘッダーからユーザー情報を取得するミドルウェアです
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// IAP は認証済みユーザーのメールアドレスをこのヘッダーで渡します
		// プレフィックス "accounts.google.com:" が付く場合があります
		email := r.Header.Get("X-Goog-Authenticated-User-Email")
		if email != "" {
			email = strings.TrimPrefix(email, "accounts.google.com:")
		}

		// コンテキストにユーザー情報を保存
		ctx := context.WithValue(r.Context(), userEmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserEmail はコンテキストからユーザーのメールアドレスを取得します
func GetUserEmail(ctx context.Context) string {
	email, ok := ctx.Value(userEmailKey).(string)
	if !ok {
		return ""
	}
	return email
}
