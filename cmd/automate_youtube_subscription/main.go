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
	google_api_key := os.Getenv("GOOGLE_API_KEY")
	target_channel_id := os.Getenv("TARGET_CHANNEL_ID")
	fmt.Printf("GOOGLE_API_KEY: %v....\nTARGET_CHANNEL_ID: %v....\n", google_api_key[:10], target_channel_id[:10])

	token, err := auth.GetValidToken(ctx)
	if err != nil {
		log.Fatalf("토큰 처리 중 오류 발생: %v", err)
	}
	client := auth.GetClient(ctx, token)

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(google_api_key))
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}
	fmt.Printf("\t======\tYOUTUBE PLAYLIST COPY START\t======\t\n")

	playListInfoMap, playListMap, err := yt.GetPlayList(youtubeService, target_channel_id)
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}
	youtubeService, err = youtube.NewService(ctx, option.WithHTTPClient(client))
	for playListTitle, playListId := range playListInfoMap {
		fmt.Printf("playListId:%v, playListTitle:%v, len(playListItems):%v\n", playListId, playListTitle, len(playListMap[playListId]))
		yt.RegisterVideoToMyPlayList(youtubeService, playListId, playListTitle, playListMap[playListId])
	}
	fmt.Printf("\t======\tYOUTUBE PLAYLIST COPY END\t======\t\n")

	fmt.Printf("\t======\tYOUTUBE SUBSCRIPTION COPY START\t======\t\n")
	target_channel_set, err := yt.GetSubscriptionSet(youtubeService, target_channel_id)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	if target_channel_set.IsEmpty() {
		fmt.Errorf("해당 채널의 구독 정보가 비공개여서 구독 목록을 가져올수 없거나 구독 목록이 없습니다.(subscribed_channel_length:%v)", target_channel_set.Size())
		return
	}
	youtubeService, err = youtube.NewService(ctx, option.WithHTTPClient(client))
	current_channel_set, err := yt.GetCurrentSubscriptionSet(youtubeService)
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}
	fmt.Printf("[willbe] retrieved subcribed channel size(%v)\n ", target_channel_set.Size())
	fmt.Printf("[current] subcribed channel size(%v) \n", current_channel_set.Size())
	target_channel_set.Subtract(current_channel_set)
	fmt.Printf("-> removed duplicated channel. left channel count: %v\n", target_channel_set.Size())
	if target_channel_set.IsEmpty() {
		fmt.Println("it's already perfectly synchronized..")
		return
	}

	for channelId := range target_channel_set.Data {
		isSuccess := yt.SubscribeToChannel(youtubeService, channelId)
		if isSuccess {
			fmt.Printf("[request] %v registered\n", channelId)
		}
	}
	fmt.Printf("\t======\tYOUTUBE SUBSCRIPTION COPY END\t======\t\n")

}
