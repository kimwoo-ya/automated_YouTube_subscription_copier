package main

import (
	"automate_youtube_subscription/internal/pkg/config"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	yt "automate_youtube_subscription/internal/pkg/youtube"
)

func main() {
	// init
	ctx := context.Background()
	config := config.GetInstance()

	if err := yt.InitializeApiKeyService(ctx); err != nil {
		log.Fatalf("Failed to init ApiKey Service %v", err)
	}
	if err := yt.InitializeOAuthService(ctx); err != nil {
		log.Fatalf("Failed to init OAuth Service %v", err)
	}
	var wg sync.WaitGroup

	// main starts..
	isPlayListCopyEnabled, _ := strconv.ParseBool(os.Getenv("PLAYLIST_COPY_ENABLED"))
	if isPlayListCopyEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			playListInfoMap, playListMap, err := yt.GetPlayList(config.GetTargetChannelId())
			if err != nil {
				log.Fatalf("%v", err)
			}
			for playListTitle, playListId := range playListInfoMap {
				fmt.Printf("playListId:%v, playListTitle:%v, len(playListItems):%v\n", playListId, playListTitle, len(playListMap[playListId]))
				yt.RegisterVideoToMyPlayList(playListId, playListTitle, playListMap[playListId])
			}
		}()
	}
	isSubscriptionListCopyEnabled, _ := strconv.ParseBool(os.Getenv("SUBSCRIPTIONLIST_COPY_ENABLED"))
	if isSubscriptionListCopyEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			target_channel_set, err := yt.GetSubscriptionSet(config.GetTargetChannelId())
			if err != nil {
				log.Fatalf("%v", err)
			}
			if target_channel_set.IsEmpty() {
				log.Fatalf("해당 채널의 구독 정보가 비공개여서 구독 목록을 가져올수 없거나 구독 목록이 없습니다.(subscribed_channel_length:%v)", target_channel_set.Size())
			}
			current_channel_set, err := yt.GetCurrentSubscriptionSet()
			if err != nil {
				log.Fatalf("%+v", err)
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
				isSuccess := yt.SubscribeToChannel(channelId)
				if isSuccess {
					fmt.Printf("[request] %v registered\n", channelId)
				}
			}
		}()
	}

	wg.Wait()
}
