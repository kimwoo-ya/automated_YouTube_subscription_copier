# youtube subscription automate copy
유튜브 프리미엄 유목민으로서.... 약 180일 마다 신규 계정을 생성해야하므로...
구 계정의 구독 목록을 신규 계정에 일괄 등록하기 위함....ㅎ
## requirements
1. google cloud registration
2. golang version `1.23.0`


## Steps
### 복제할 유튜브 계정의 채널 아이디 정보를 `.env`파일에 기입해주세요.
```bash
touch .env
```
```.env
TARGET_CHANNEL_ID="PASTE_TARGET_CHANNEL_ID"
```

1. 채널의 구독 정보 비공개를 해제해주세요.<br/>
![ ](./screenshots/00_사전조치사항.png)
2. 구글 클라우드 콘솔로 이동해주세요. [google-cloud-console](https://console.cloud.google.com/welcome?hl=ko&inv=1&invt=Ab0cDg)
3. 신규 프로젝트 생성 ![](./screenshots/01_리소스%20생성.png)
4. 신규로 리소스를 생성해주세요. ![](./screenshots/01-1.png)
5. Oauth2 client를 생성해주세요.(구독 요청시 이용)
![](./screenshots/02_0oauth%20클라이언트%20만들기.png)
![](./screenshots/02-1.png)
![](./screenshots/02-2.png)
![](./screenshots/02-3.png)
![](./screenshots/02-4.png)
![](./screenshots/02-5.png)
중요 개인 정보를 `.env`파일에 기입해주세요.
```.env
GOOGLE_CLIENT_ID="PASTE_YOUR_CLIENT_ID"
GOOGLE_CLIENT_SECRET="PASTE_YOUR_CLIENT_SECRET"
REDIRECT_URL="http://localhost:8080"
```
6. API Key를 생성해주세요.(구독 목록 조회시 이용)
![](./screenshots/03-0APIKEY만들기.png)
![](./screenshots/03-1.png)
![](./screenshots/03-2.png)
save private informations to file `.env`
```.env
GOOGLE_API_KEY="PASTE_YOUR_API_KEY"
```
6-1. youtube data api v3 사용 등록
![](./screenshots/03-3.png)
![](./screenshots/03-4.png)
![](./screenshots/03-5.png)
![](./screenshots/03-6.png)
![](./screenshots/03-7.png)
7. 테스트 용 사용자 추가(oauth 인증시 필요)
![](./screenshots/04_앱게시.png)
![](./screenshots/04-2.png)
8. 명령어 실행
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