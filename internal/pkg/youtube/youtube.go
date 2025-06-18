package youtube

import (
	"fmt"
	"log"

	"google.golang.org/api/youtube/v3"
)

func GetSubscriptionList(youtubeService *youtube.Service, target_channel_id string) ([]string, error) {
	channel_list := []string{}
	call := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"})
	call = call.ChannelId(target_channel_id)
	call = call.MaxResults(100000)
	response, err := call.Do()
	if err != nil {
		return channel_list, fmt.Errorf("%v", err)
	}
	fmt.Printf("\n==\t채널아이디(%v)의 YouTube 구독 목록\t==\n", target_channel_id)
	for _, item := range response.Items {
		fmt.Printf("- 채널 제목: %s (ID: %s)\n", item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
		channel_list = append(channel_list, item.Snippet.ResourceId.ChannelId)
	}

	// 다음 페이지가 있다면 토큰을 사용하여 다음 페이지를 가져올 수 있습니다.
	for response.NextPageToken != "" {
		nextCall := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"})
		nextCall = nextCall.ChannelId(target_channel_id)
		nextCall = nextCall.MaxResults(100000)
		nextCall = nextCall.PageToken(response.NextPageToken)
		nextResponse, nextErr := nextCall.Do()
		if nextErr != nil {
			log.Fatalf("다음 페이지 구독 목록 가져오기 오류: %v", nextErr)
			return []string{}, fmt.Errorf("%v", nextErr)
		}
		for _, item := range nextResponse.Items {
			fmt.Printf("- 채널 제목: %s (ID: %s)\n", item.Snippet.Title, item.Snippet.ResourceId.ChannelId)
			channel_list = append(channel_list, item.Snippet.ResourceId.ChannelId)
		}
		response = nextResponse
	}
	return channel_list, nil
}

func SubscribeToChannel(service *youtube.Service, channelId string) {
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
	}

	fmt.Printf("구독 성공: %s\n", response.Snippet.Title)
}
