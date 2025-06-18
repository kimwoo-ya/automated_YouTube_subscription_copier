# youtube subscription automate copy

## requirements
1. google cloud registration
2. golang version `1.23.0`


## Steps
### save target channel ID informations to file `.env`
```bash
touch .env
```
```.env
TARGET_CHANNEL_ID="PASTE_TARGET_CHANNEL_ID"
```

1. Please set the subscription list to public so it can be duplicated.![ ](./screenshots/00_사전조치사항.png)
2. Move to [google-cloud-console](https://console.cloud.google.com/welcome?hl=ko&inv=1&invt=Ab0cDg)
3. Create New project ![](./screenshots/01_리소스%20생성.png)
4. Register New resource ![](./screenshots/01-1.png)
5. Make Oauth2 client
![](./screenshots/02_0oauth%20클라이언트%20만들기.png)
![](./screenshots/02-1.png)
![](./screenshots/02-2.png)
![](./screenshots/02-3.png)
![](./screenshots/02-4.png)
![](./screenshots/02-5.png)
save private informations to file `.env`
```.env
GOOGLE_CLIENT_ID="PASTE_YOUR_CLIENT_ID"
GOOGLE_CLIENT_SECRET="PASTE_YOUR_CLIENT_SECRET"
REDIRECT_URL="http://localhost:8080"
```
6. Make API Key
![](./screenshots/03-0APIKEY만들기.png)
![](./screenshots/03-1.png)
![](./screenshots/03-2.png)
save private informations to file `.env`
```.env
GOOGLE_API_KEY="PASTE_YOUR_API_KEY"
```
![](./screenshots/03-3.png)
![](./screenshots/03-4.png)
![](./screenshots/03-5.png)
![](./screenshots/03-6.png)
![](./screenshots/03-7.png)
7. publish app
![](./screenshots/04_앱게시.png)
![](./screenshots/04-2.png)
8. execute code.
```bash
$ go run cmd/automate_youtube_subscription/main.go
GOOGLE_CLIENT_ID: XXXX....
GOOGLE_CLIENT_SECRET: XXXX....
REDIRECT_URL: http://localhost:8080
GOOGLE_API_KEY: XXXX....
TARGET_CHANNEL_ID: XXXX....

==	채널아이디의 YouTube 구독 목록	==
- 채널 제목: Noel Deyzel (ID: UCMp-0bU-PA7BNNR-zIvEydA)
......
```