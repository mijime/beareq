![](docs/assets/logo.png)

## Configuration

`~/.config/beareq/profiles.toml`

```toml
[google]
ClientID = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com"
ClientSecret = "xxxxxxxxxxxxxxxxxxxxxxxx"
RedirectURL = "http://localhost:8999"
Scopes = ["https://www.googleapis.com/auth/tasks"]

[google.Endpoint]
AuthStyle = 2
AuthURL = "https://accounts.google.com/o/oauth2/auth"
TokenURL = "https://www.googleapis.com/oauth2/v3/token"

[slack]
ClientID = "000000000000.0000000000000"
ClientSecret = "00000000000000000000000000000000"
RedirectURL = "http://localhost:8999"
Scopes = ["chat:write:user", "chat:write:bot"]

[slack.Endpoint]
AuthURL = "https://slack.com/oauth/authorize"
TokenURL = "https://slack.com/api/oauth.access"
```

`~/.config/beareq/commands.toml`

```toml
[slack.chatPostMessage]

URL = "https://slack.com/api/chat.postMessage"
Method = "POST"
Data = """
{
  "channel": "{{ index .channel 0 }}",
  "text": "{{ index .text 0 }}"
}
"""

[slack.chatPostMessage.Header]
Content-Type = ["application/json; charset=utf8"]
```

## Usage

```bash
beareq --profile slack --header=Content-type:\ application/json \
  --data='{"channel":"CXXXXX","text":"Helloworld"}' https://slack.com/api/chat.postMessage
```

```bash
beareq --profile slack "slack://chatPostMessage?channel=CXXXXX&text=:smile:"
```
