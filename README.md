# Telelog

Embed to app or tail a file, the log will send to your telegram (chat_id)

## Installation

Using Go Modules

```bash
go get github.com/codenoid/telelog
```

## Usage

Example of recipient file : 

```txt
999838460
961268461
827957192
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

    err := logger.Connect()
    if err != nil {
        panic(err)
    }

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