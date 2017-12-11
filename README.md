# pigeon
> Notify your own events on telegram using RPC requests and command line.

Pigeon is a `server` `client` program that can send message using command line to telegram using a telegram bot.

### Installation

```
    go get -u github.com/ahmdrz/pigeon/...
```

### Server

First of all start server.

```
    Usage of pigeon-server:
        -output string
                Output of pigeon.
        -port int
                Listen port (default 8765)
        -token string
                Telegram bot token
        -users string
                Telegram ID of users. separate it by ','

    Notice: To find your telegram id, you can use @userinfobot
    Example: pigeon-server -port=8595 -token='<BOT-TOKEN>' -users='<USER1>,<USER2>'
```

### Client

```
    Usage of pigeon-client:
        -body string
                body|content of message
        -port int
                pigeon server port (default 8765)
        -server string
                pigeon server address (default "localhost")
        -title string
                title of message
```