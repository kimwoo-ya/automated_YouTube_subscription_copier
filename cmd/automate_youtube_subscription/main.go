package main

import (
	"automate_youtube_subscription/internal/pkg/auth"
	"context"
	"fmt"
	"log"
	"os"

	yt "automate_youtube_subscription/internal/pkg/youtube"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	token, err := auth.GetValidToken(ctx)
	if err != nil {
		log.Fatalf("토큰 처리 중 오류 발생: %v", err)
	}
	client := auth.GetClient(ctx, token)

	google_api_key := os.Getenv("GOOGLE_API_KEY")
	target_channel_id := os.Getenv("TARGET_CHANNEL_ID")
	fmt.Printf("GOOGLE_API_KEY: %v....\nTARGET_CHANNEL_ID: %v....\n", google_api_key[:10], target_channel_id[:10])
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(google_api_key))
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}
	channel_list, err := yt.GetSubscriptionList(youtubeService, target_channel_id)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	if len(channel_list) == 0 {
		fmt.Errorf("해당 채널의 구독 정보가 비공개여서 구독 목록을 가져올수 없거나 구독 목록이 없습니다.(subscribed_channel_length:%v)", len(channel_list))
		return
	}

	youtubeService, err = youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}
	for _, channel := range channel_list {
		yt.SubscribeToChannel(youtubeService, channel)
	}
}
