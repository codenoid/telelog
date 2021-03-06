# Telelog

[![Go Reference](https://pkg.go.dev/badge/github.com/codenoid/telelog.svg)](https://pkg.go.dev/github.com/codenoid/telelog)

Embed to app or tail a file, the log will send to your telegram (chat_id)

![Telegram Notification](telelog.jpg?raw=true)

## Installation

Embed to your app using Go Modules

```bash
go get github.com/codenoid/telelog
```

CLI version

```bash
go get github.com/codenoid/telelog/cmd/telelog
```

## Usage

Example of recipient file content (text file) : 

```txt
999838460
961268461
827957192
```

CMD Example :

you must set env var for cli version !

```bash
$ codenoid> telelog 
Telelog, make sure you already set these env : 
TELELOG_BOT_TOKEN
TELELOG_APP_NAME
TELELOG_DEBUG_MODE (optional)
TELELOG_RECIPIENT_LIST

$ codenoid> telelog -file ./path/to/error.log -level error
```

Code Example :

```go
package main

import (
    "time"

    "github.com/codenoid/telelog"
)

func main() {
    logger := telelog.LoggerNew()

    // default TELELOG_BOT_TOKEN, unless you call SetToken
    logger.SetToken("1125121251:AAF2sfBCbKjag8LhUIAzf1mzk36BxcJ0Mvg")

    // default TELELOG_APP_NAME, unless you call SetAppName
    logger.SetAppName("Uploader Service")

    // default TELELOG_DEBUG_MODE, unless you call SetDebug
    // if true, any logger.Debug* log will be send 
    logger.SetDebug(false)

    // default TELELOG_RECIPIENT_LIST (single path to file), unless you call SetRecipient
    logger.SetRecipient("/path/to/text/file.txt", "/second/file/that/contain/chat_id.txt")

    // or reader
    f, _ := os.Open("file.txt")
    logger.SetRecipientFromReader(f)

    // or []byte
    logger.SetRecipientFromByte([]byte{"777000\n"})

    logger.Warn("Warning! your app will be error")
    logger.Error("yo, this is error, in 2sec your app will dead")
    time.Sleep(2*time.Second)
    logger.Fatal("lmaoo")

}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
