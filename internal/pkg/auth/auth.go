package auth

import (
	"automate_youtube_subscription/internal/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

const tokenFile = "token.json"

var (
	codeCh = make(chan string)
)

func loadToken() (*oauth2.Token, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	if token.Expiry.Before(time.Now()) {
		return &token, fmt.Errorf("토큰 만료됨")
	}
	return &token, nil
}

func GetValidToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := loadToken()
	if err != nil || !token.Valid() {
		fmt.Println("토큰 없음 또는 만료됨. 인증 진행.")
		oauth_config := config.GetInstance().GetOauthConfig()
		authURL := oauth_config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		fmt.Println("브라우저에서 URL 열기:", authURL)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "인증 코드가 없습니다", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "인증 성공! 터미널로 돌아가세요.")
			codeCh <- code
		})
		go func() {
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()

		code := <-codeCh
		token, err = oauth_config.Exchange(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("토큰 교환 실패: %w", err)
		}
		saveToken(token)
	}
	return token, nil
}

func GetClient(ctx context.Context, token *oauth2.Token) *http.Client {
	oauth_config := config.GetInstance().GetOauthConfig()
	return oauth_config.Client(ctx, token)
}

func saveToken(token *oauth2.Token) {
	data, _ := json.Marshal(token)
	_ = os.WriteFile(tokenFile, data, 0600)
}
