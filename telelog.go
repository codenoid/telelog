package telelog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/withmandala/go-log"
)

const (
	FATAL = "💥 [FATAL]"
	ERROR = "🥴 [ERROR]"
	WARN  = "⚠️ [WARN]"
	INFO  = "👀 [INFO]"
	DEBUG = "🤔 [DEBUG]"
)

// Logger telelog logger instance
type Logger struct {
	name      string
	debug     bool
	token     string
	file      string
	recipient []int64

	osLogger *log.Logger

	bot *tgbotapi.BotAPI
}

// LoggerNew create new telelog Instance
func LoggerNew(token string) *Logger {

	i := &Logger{}

	i.osLogger = log.New(os.Stderr).WithColor()

	if token := os.Getenv("TELELOG_BOT_TOKEN"); token != "" {
		i.token = token
	}

	if name := os.Getenv("TELELOG_APP_NAME"); name != "" {
		i.name = name
	}

	if debug := os.Getenv("TELELOG_DEBUG_MODE"); debug != "" {
		switch strings.ToLower(debug) {
		case "1", "true", "enabled", "active", "yes":
			i.debug = true
		}
	}

	if path := os.Getenv("TELELOG_RECIPIENT_LIST"); path != "" {
		i.SetRecipientFromFiles(path)
	}

	// set and create telegram instance
	i.SetToken(token)

	return i
}

func (i *Logger) SetToken(token string) error {
	token = strings.TrimSpace(token)

	if token != i.token {
		i.token = token

		bot, err := tgbotapi.NewBotAPI(i.token)
		if err != nil {
			return err
		}

		i.bot = bot
	}

	return nil
}

func (i *Logger) SetAppName(name string) {
	i.name = name
}

func (i *Logger) SetDebug(debug bool) {
	i.debug = debug
}

func (i *Logger) SetRecipientFromFiles(files ...string) {
	// iterate given list of file path
	for _, path := range files {
		// read file lines as slice of string
		if lines, err := readLines(path); err == nil {
			i.setRecipient(lines)
		}
	}
}

// SetRecipientFromByte receive byte of file content
func (i *Logger) SetRecipientFromByte(b []byte) {
	i.setRecipient(
		stringSplitLines(
			string(b),
		),
	)
}

// SetRecipientFromReader read content from reader
func (i *Logger) SetRecipientFromReader(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	i.setRecipient(
		stringSplitLines(
			string(b),
		),
	)

	return nil
}

func (i *Logger) setRecipient(recipient []string) {
	// iterate file that contain string of chat_id
	for _, chatIDStr := range recipient {
		// convert string chat_id into int64
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			i.osLogger.Warn(err.Error())
			continue
		}
		if chatID > 0 {
			i.recipient = append(i.recipient, chatID)
		}
	}
}

func (i *Logger) sendLog(level string, msg string) {

	content := `%v %v

Filename: %v
Line: %v
FuncName: %v

Message:
%v`

	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		file, line := details.FileLine(pc)
		fileName := file
		// privacy first
		if val := strings.Split(file, string(os.PathSeparator)); len(val) > 0 {
			fileName = val[len(val)-1]
		}
		content = fmt.Sprintf(content, level, i.name, fileName, line, details.Name(), msg)
	}

	go func() {
		for _, chatID := range i.recipient {
			msg := tgbotapi.NewMessage(chatID, content)
			i.bot.Send(msg)
		}
	}()

}

// Fatal print fatal message to output and quit the application with status 1
func (i *Logger) Fatal(v ...interface{}) {
	i.sendLog(FATAL, fmt.Sprintln(v...))
	panic(v)
}

// Fatalf print formatted fatal message to output and quit the application
// with status 1
func (i *Logger) Fatalf(format string, v ...interface{}) {
	i.sendLog(FATAL, fmt.Sprintf(format, v...))
	panic(v)
}

// Error print error message to output
func (i *Logger) Error(v ...interface{}) {
	i.sendLog(ERROR, fmt.Sprintln(v...))
}

// Errorf print formatted error message to output
func (i *Logger) Errorf(format string, v ...interface{}) {
	i.sendLog(ERROR, fmt.Sprintf(format, v...))
}

// Warn print warning message to output
func (i *Logger) Warn(v ...interface{}) {
	i.sendLog(WARN, fmt.Sprintln(v...))
}

// Warnf print formatted warning message to output
func (i *Logger) Warnf(format string, v ...interface{}) {
	i.sendLog(WARN, fmt.Sprintf(format, v...))
}

// Info print informational message to output
func (i *Logger) Info(v ...interface{}) {
	i.sendLog(INFO, fmt.Sprintln(v...))
}

// Infof print formatted informational message to output
func (i *Logger) Infof(format string, v ...interface{}) {
	i.sendLog(INFO, fmt.Sprintf(format, v...))
}

// Debug print debug message to output if debug output enabled
func (i *Logger) Debug(v ...interface{}) {
	if i.debug {
		i.sendLog(DEBUG, fmt.Sprintln(v...))
	}
}

// Debugf print formatted debug message to output if debug output enabled
func (i *Logger) Debugf(format string, v ...interface{}) {
	if i.debug {
		i.sendLog(DEBUG, fmt.Sprintf(format, v...))
	}
}
