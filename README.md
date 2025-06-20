# Automated YouTube Subscription Copier
> As a YouTube Premium nomad, <br/> I need to create/transfer a new account every 180 days or so, <br/>
To register the subscription list of the old account to the new account in bulk. <br/>
There are similar functions in open source, but they are not properly maintained or are performed using `selenium`, `Puppeteer`<br/>
→ They don't work well or takes time too long. 😥

#### `The current functionality appears to be working normally, but once we determine that it is fully functional, we will release it as a major version.`

## TODO
1. Add test cases
2. Do I really need to query the subscription list using OPENAPI?<br/> [I think it will work if I parse it well on this page....](https://www.youtube.com/feed/channels)



### Key Features
1. Copy YouTube Subscriptions (Original Account → New Account)
2. Copy YouTube Playlists (ex. YouTube Music) (Original Account → New Account)


## requirements
1. google cloud registration
2. Choose one of the two
    - golang
    ```bash
        $ go version
        go version go1.23.0 darwin/arm64
    ```
    - docker
    ```bash
        $ docker version
        Server: Docker Desktop 4.26.1 (131620)
        Engine:
        Version:          24.0.7
    ```
3. Obtain original channel ID (old account)


## Steps
Please enter the channel ID information of the YouTube account you want to clone into the `.env` file.<br/>
0. Log in with your original channel ID (old account)
- [this link ](https://www.youtube.com/account_advanced) provides your channel id informations
![](./screenshots/채널%20아이디%20조회.png)
```bash
touch .env
# append below
TARGET_CHANNEL_ID="PASTE_TARGET_CHANNEL_ID"
```
1. Set publically the subscription information of the original channel (old account).<br/>
![ ](./screenshots/00_사전조치사항.png)
2. Please go to the Google Cloud Console. <br/>[google-cloud-console](https://console.cloud.google.com/welcome?hl=ko&inv=1&invt=Ab0cDg)
3. Create a new project ![](./screenshots/01_리소스%20생성.png)
4. Please create a new resource. ![](./screenshots/01-1.png)
5. Please create an Oauth2 client (to be used when requesting a subscription)
![](./screenshots/02_0oauth%20클라이언트%20만들기.png)
![](./screenshots/02-1.png)
![](./screenshots/02-2.png)
![](./screenshots/02-3.png)
![](./screenshots/02-4.png)
![](./screenshots/02-5.png)
6. Please enter obtained information in the `.env` file.
```.env
GOOGLE_CLIENT_ID="PASTE_YOUR_CLIENT_ID"
GOOGLE_CLIENT_SECRET="PASTE_YOUR_CLIENT_SECRET"
REDIRECT_URL="http://localhost:8080"
```
7. Please create an API Key (used when viewing subscription list)
![](./screenshots/03-0APIKEY만들기.png)
![](./screenshots/03-1.png)
![](./screenshots/03-2.png)
save private informations to file `.env`
```.env
GOOGLE_API_KEY="PASTE_YOUR_API_KEY"
```

8. Register to use youtube data api v3
![](./screenshots/03-3.png)
![](./screenshots/03-4.png)
![](./screenshots/03-5.png)
![](./screenshots/03-6.png)
![](./screenshots/03-7.png)

9. Add a user for testing (required for oauth authentication)
![](./screenshots/04_앱게시.png)
![](./screenshots/04-2.png)

10. Everything is ready. Now, shall we give it a try?
- Running with Docker
```bash
# build docker images from dockerfile
# Here is an example on arm mac.
$ docker buildx build --platform linux/amd64 -t automate_youtube_subscription -f internal/deployments/Dockerfile .

# execute container.
$ docker run --rm -p 8080:8080 --name automate_youtube_subscription automate_youtube_subscription
```
- Running with Go
```bash
$ go run cmd/automate_youtube_subscription/main.go
```
### output
```bash
GOOGLE_CLIENT_ID: XXXX....
GOOGLE_CLIENT_SECRET: XXXX....
REDIRECT_URL: http://localhost:8080
GOOGLE_API_KEY: XXXX....
TARGET_CHANNEL_ID: XXXX....

- 재생 목록: 🔥 (ID:XXXXXX)
- 재생 목록 항목: Make U Dance (Feat. Jay Park) (박재범) (& Paul Blanco) (ID:4DZRLuD8AMs)
- 재생 목록 항목: Crack On My Screen (내핸드폰에금이갔네) (Prod. By Minit) (Feat. Paloalto) (ID:kgNiM2u9OHQ)
- 재생 목록 항목: [MV] Just Music _ Carnival Gang(카니발갱) (ID:CJdOUxMAkME)
....
playListId:XYYYY, playListTitle:🔥, len(playListItems):10
playlist is already exist. 복제된_🔥
이미 있는 재생목록 아이템(4DZRLuD8AMs)
이미 있는 재생목록 아이템(kgNiM2u9OHQ)
.....
playlist item is newly appended [MV] Just Music _ Carnival Gang(카니발갱)
.....

==	채널아이디(XXX)의 YouTube 구독 목록	==
- 채널 제목: Noel Deyzel (ID: UCMp-0bU-PA7BNNR-zIvEydA)
.....

==	내 채널의 YouTube 구독 목록	==
- 채널 제목: acooknamedMatt (ID: UCYjJeNVpgjAz-Sv4dhs13VQ)
.....

[willbe] retrieved subcribed channel size(85)
[current] subcribed channel size(84)
→ removed duplicated channel. left channel count:
구독 성공: Noel Deyzel
[request] UCMp-0bU-PA7BNNR-zIvEydA registered
%
......
```

## reference
- [OPEN API usage limits](https://developers.google.com/youtube/v3/determine_quota_cost)
- [Use case for YouTube Music Playlists](https://developers.google.com/youtube/v3/docs/playlists)
- [Official Docs](https://developers.google.com/youtube/v3/quickstart/go#step_1_turn_on_the)
