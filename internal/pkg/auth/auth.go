package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const tokenFile = "token.json"

var (
	conf   *oauth2.Config
	codeCh = make(chan string)
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	fmt.Printf("GOOGLE_CLIENT_ID: %v....\nGOOGLE_CLIENT_SECRET: %v....\nREDIRECT_URL: %v\n", os.Getenv("GOOGLE_CLIENT_ID")[:10], os.Getenv("GOOGLE_CLIENT_SECRET")[:10], os.Getenv("REDIRECT_URL"))
	conf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/youtube"},
		Endpoint:     google.Endpoint,
	}
}

func GetValidToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := loadToken()
	if err != nil || !token.Valid() {
		fmt.Println("토큰 없음 또는 만료됨. 인증 진행.")
		authURL := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		fmt.Println("브라우저에서 URL 열기:", authURL)

		http.HandleFunc("/", handleOAuthCallback)
		go func() {
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()

		code := <-codeCh
		token, err = conf.Exchange(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("토큰 교환 실패: %w", err)
		}
		saveToken(token)
	}
	return token, nil
}

func GetClient(ctx context.Context, token *oauth2.Token) *http.Client {
	return conf.Client(ctx, token)
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "인증 코드가 없습니다", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "인증 성공! 터미널로 돌아가세요.")
	codeCh <- code
}

func saveToken(token *oauth2.Token) {
	data, _ := json.Marshal(token)
	_ = os.WriteFile(tokenFile, data, 0600)
}

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
