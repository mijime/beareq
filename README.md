![](docs/assets/logo.png)

## Configuration

`~/.config/beareq/profiles.toml`

```toml
[google.OAuth]
ClientID = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com"
ClientSecret = "xxxxxxxxxxxxxxxxxxxxxxxx"
RedirectURL = "http://localhost:8999"
Scopes = ["https://www.googleapis.com/auth/tasks"]

[google.OAuth.Endpoint]
AuthStyle = 2
AuthURL = "https://accounts.google.com/o/oauth2/auth"
TokenURL = "https://www.googleapis.com/oauth2/v3/token"

[slack.OAuth]
ClientID = "000000000000.0000000000000"
ClientSecret = "00000000000000000000000000000000"
RedirectURL = "https://localhost:8999"
Scopes = ["chat:write:user", "chat:write:bot"]

[slack.OAuth.Endpoint]
AuthURL = "https://slack.com/oauth/authorize"
TokenURL = "https://slack.com/api/oauth.access"

[[slack.OpenAPI]] # OPTIONAL
BaseURL = "https://api.slack.com/api" # OPTIONAL
Specs = ["https://raw.githubusercontent.com/slackapi/slack-api-specs/master/web-api/slack_web_openapi_v2_without_examples.json"]

[notion]

[notion.Header]
Notion-Version = ["2021-05-13"]
Authorization = ["Bearer secret_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"]
```

## Usage

```bash
beareq \
	--profile slack \
	--header=Content-type:\ application/json \
	--data='{"channel":"CXXXXX","text":"Helloworld"}' \
	https://slack.com/api/chat.postMessage
```

```bash
beareq \
	--profile slack \
	--header=Content-type:\ application/json \
	--jo channel=CXXXXX \
	--jo text=Helloworld \
	https://slack.com/api/chat.postMessage
```

```bash
beareq-oapi \
	--profile slack \
	chat_postMessage \
	--channel CXXXXX \
	--text "Helloworld"
```
