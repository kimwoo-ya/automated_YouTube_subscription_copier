package youtube

import (
	"fmt"
	"log"
	"strings"

	datatype "automate_youtube_subscription/internal/pkg/utils/data_type"

	"google.golang.org/api/youtube/v3"
)

var PAGE_SIZE int64

func init() {
	PAGE_SIZE = 30
}

func GetPlayList(youtubeService *youtube.Service, target_channel_id string) (map[string]string, map[string][]string, error) {
	playListInfo := make(map[string]string)
	playList := make(map[string][]string, 0)

	call := youtubeService.Playlists.List([]string{"snippet", "contentDetails"}).ChannelId(target_channel_id)
	resp, err := call.Do()
	if err != nil {
		fmt.Errorf("%+v", err)
		return playListInfo, playList, err
	}

	for _, item := range resp.Items {
		playListInfo[item.Snippet.Title] = item.Id
		playList[item.Id] = make([]string, 0)
	}
	for playListTitle, playListId := range playListInfo {
		fmt.Printf("- 재생 목록: %s (ID:%v)\n", playListTitle, playListId)

		innerCall := youtubeService.PlaylistItems.List([]string{"snippet", "contentDetails"}).PlaylistId(playListId).MaxResults(PAGE_SIZE)
		response, err := innerCall.Do()
		if err != nil {
			fmt.Errorf("%+v", err)
			return playListInfo, playList, err
		}

		for _, playlistItem := range response.Items {
			fmt.Printf("- 재생 목록 항목: %s (ID:%v)\n", playlistItem.Snippet.Title, playlistItem.Snippet.ResourceId.VideoId)
			playList[playListId] = append(playList[playListId], playlistItem.Snippet.ResourceId.VideoId)
		}

		for response.NextPageToken != "" {
			nextCall := youtubeService.PlaylistItems.List([]string{"snippet", "contentDetails"}).PlaylistId(playListId).MaxResults(PAGE_SIZE).PageToken(response.NextPageToken)
			nextResponse, nextErr := nextCall.Do()
			if nextErr != nil {
				log.Fatalf("다음 페이지 구독 목록 가져오기 오류: %v", nextErr)
				return nil, nil, fmt.Errorf("%v", nextErr)
			}
			for _, playlistItem := range nextResponse.Items {
				fmt.Printf("- 재생 목록 항목: %s (ID:%v)\n", playlistItem.Snippet.Title, playlistItem.Snippet.ResourceId.VideoId)
				playList[playListId] = append(playList[playListId], playlistItem.Snippet.ResourceId.VideoId)
			}
			response = nextResponse

		}

	}

	return playListInfo, playList, nil
}

func RegisterVideoToMyPlayList(youtubeService *youtube.Service, playListId string, playListTitle string, videoIdList []string) {
	targetPlaylistTitle := fmt.Sprintf("복제된_%v", playListTitle)
	existingPlaylist, err := findPlaylistByTitle(youtubeService, targetPlaylistTitle)
	if err != nil {
		log.Fatalf("재생목록을 찾는 중 오류 발생: %v", err)
		return
	}
	var createdPlaylist *youtube.Playlist
	existingVideoIDs := make(map[string]bool, 0)
	if existingPlaylist == nil {
		myPlaylist := &youtube.Playlist{
			Snippet: &youtube.PlaylistSnippet{
				Title:       targetPlaylistTitle,
				Description: fmt.Sprintf("원본: %v,", playListId),
			},
			Status: &youtube.PlaylistStatus{PrivacyStatus: "private"},
		}
		call := youtubeService.Playlists.Insert([]string{"snippet", "status"}, myPlaylist)
		createdPlaylist, err = call.Do()
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
		fmt.Printf("playlist is newly generated %v\n", createdPlaylist.Snippet.Title)
	} else {
		createdPlaylist = existingPlaylist
		fmt.Printf("playlist is already exist. %v\n", createdPlaylist.Snippet.Title)
		pageToken := ""
		for {
			call := youtubeService.PlaylistItems.List([]string{"snippet"}).PlaylistId(createdPlaylist.Id).MaxResults(PAGE_SIZE).PageToken(pageToken)
			response, err := call.Do()
			if err != nil {
				log.Fatalf("기 존재 재생목록 아이템을 가져오는 중 오류 발생: %v", err)
				return
			}
			for _, item := range response.Items {
				if item.Snippet.ResourceId.Kind == "youtube#video" {
					existingVideoIDs[item.Snippet.ResourceId.VideoId] = true
				}
			}
			pageToken = response.NextPageToken
			if pageToken == "" {
				break
			}
		}
	}

	for _, videoId := range videoIdList {
		if existingVideoIDs[videoId] {
			fmt.Printf("이미 있는 재생목록 아이템(%v)\n", videoId)
			continue
		}
		playlistItem := &youtube.PlaylistItem{
			Snippet: &youtube.PlaylistItemSnippet{
				PlaylistId: createdPlaylist.Id,
				ResourceId: &youtube.ResourceId{
					Kind:    "youtube#video",
					VideoId: videoId,
				},
			},
		}
		res, err := youtubeService.PlaylistItems.Insert([]string{"snippet"}, playlistItem).Do()
		if err != nil {
			log.Fatal("%v", err)
			continue
		}
		fmt.Printf("playlist item is newly appended %v\n", res.Snippet.Title)

	}
}

func findPlaylistByTitle(youtubeService *youtube.Service, playListTitle string) (*youtube.Playlist, error) {
	call := youtubeService.Playlists.List([]string{"snippet"}).Mine(true).MaxResults(50)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("재생목록을 검색할 수 없습니다: %v", err)
	}

	for _, playlist := range response.Items {
		if strings.Contains(playlist.Snippet.Title, playListTitle) {
			return playlist, nil
		}
	}
	return nil, nil
}

func GetSubscriptionSet(youtubeService *youtube.Service, target_channel_id string) (*datatype.Set[string], error) {
	subscription_set := datatype.NewSet[string]()
	call := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).ChannelId(target_channel_id).MaxResults(PAGE_SIZE)
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
		nextCall := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).ChannelId(target_channel_id).MaxResults(PAGE_SIZE).PageToken(response.NextPageToken)
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
	call := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(PAGE_SIZE)
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
		nextCall := youtubeService.Subscriptions.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(PAGE_SIZE).PageToken(response.NextPageToken)
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
