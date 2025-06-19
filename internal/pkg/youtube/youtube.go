package youtube

import (
	"fmt"
	"log"

	datatype "automate_youtube_subscription/internal/pkg/utils/data_type"

	"google.golang.org/api/youtube/v3"
)

func GetSubscriptionSet(youtubeService *youtube.Service, target_channel_id string) (*datatype.Set[string], error) {
	subscription_set := datatype.NewSet[string]()
	call := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).ChannelId(target_channel_id).MaxResults(100000)
	response, err := call.Do()
	if err != nil {
		return subscription_set, fmt.Errorf("%v", err)
	}
	fmt.Printf("\n==\t채널아이디(%v)의 YouTube 구독 목록\t==\n", target_channel_id)
	for _, item := range response.Items {
		fmt.Printf("- 채널 제목: %s (ID: %s)\n", item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
		subscription_set.Add(item.Snippet.ResourceId.ChannelId)
	}

	for response.NextPageToken != "" {
		nextCall := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).ChannelId(target_channel_id).MaxResults(100000).PageToken(response.NextPageToken)
		nextResponse, nextErr := nextCall.Do()
		if nextErr != nil {
			log.Fatalf("다음 페이지 구독 목록 가져오기 오류: %v", nextErr)
			return datatype.NewSet[string](), fmt.Errorf("%v", nextErr)
		}
		for _, item := range nextResponse.Items {
			fmt.Printf("- 채널 제목: %s (ID: %s)\n", item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
			subscription_set.Add(item.Snippet.ResourceId.ChannelId)
		}
		response = nextResponse
	}
	return subscription_set, nil
}

func GetCurrentSubscriptionSet(youtubeService *youtube.Service) (*datatype.Set[string], error) {
	subscription_set := datatype.NewSet[string]()
	call := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(10000)
	response, err := call.Do()
	if err != nil {
		return subscription_set, fmt.Errorf("%v", err)
	}
	fmt.Printf("\n==\t내 채널의 YouTube 구독 목록\t==\n")
	for _, item := range response.Items {
		fmt.Printf("- 채널 제목: %s (ID: %s)\n", item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
		subscription_set.Add(item.Snippet.ResourceId.ChannelId)
	}

	for response.NextPageToken != "" {
		nextCall := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(100000).PageToken(response.NextPageToken)
		nextResponse, nextErr := nextCall.Do()
		if nextErr != nil {
			log.Fatalf("다음 페이지 구독 목록 가져오기 오류: %v", nextErr)
			return datatype.NewSet[string](), fmt.Errorf("%v", nextErr)
		}
		for _, item := range nextResponse.Items {
			fmt.Printf("- 채널 제목: %s (ID: %s)\n", item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
			subscription_set.Add(item.Snippet.ResourceId.ChannelId)
		}
		response = nextResponse
	}
	return subscription_set, nil
}

func SubscribeToChannel(service *youtube.Service, channelId string) bool {
	sub := &youtube.Subscription{
		Snippet: &youtube.SubscriptionSnippet{
			ResourceId: &youtube.ResourceId{
				Kind:      "youtube#channel",
				ChannelId: channelId,
			},
		},
	}

	call := service.Subscriptions.Insert([]string{"snippet"}, sub)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("구독 요청 실패: %v", err)
		return false
	}

	fmt.Printf("구독 성공: %s\n", response.Snippet.Title)
	return true
}
