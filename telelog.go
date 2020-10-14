package telelog

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

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

type instance struct {
	name      string
	debug     bool
	token     string
	file      string
	recipient []int64

	osLogger *log.Logger

	bot *tgbotapi.BotAPI
}

func (i *instance) LoggerNew() *instance {

	i.osLogger = log.New(os.Stderr).WithColor()

	if path := os.Getenv("TELELOG_RECIPIENT_LIST"); path != "" {
		// read file lines as slice of string
		if lines, err := readLines(path); err == nil {
			// iterate file that contain string of chat_id
			for _, chatIDStr := range lines {
				// convert string chat_id into int64
				chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
				if err != nil {
					i.osLogger.Warn(err.Error())
					continue
				}
				i.recipient = append(i.recipient, chatID)
			}
		}
	}

	return i
}

func (i *instance) SetToken(token string) {
	i.token = token
}

func (i *instance) SetAppName(name string) {
	i.name = name
}

func (i *instance) SetDebug(debug bool) {
	i.debug = debug
}

func (i *instance) SetRecipient(files ...string) {
	// iterate given list of file path
	for _, path := range files {
		// read file lines as slice of string
		if lines, err := readLines(path); err == nil {
			// iterate file that contain string of chat_id
			for _, chatIDStr := range lines {
				// convert string chat_id into int64
				chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
				if err != nil {
					i.osLogger.Warn(err.Error())
					continue
				}
				i.recipient = append(i.recipient, chatID)
			}
		}
	}
}

func (i *instance) Connect() error {
	bot, err := tgbotapi.NewBotAPI(i.token)
	if err != nil {
		return err
	}

	i.bot = bot

	return nil
}

func (i *instance) SendLog(level string, msg interface{}) {

	content := `%v

Filename: %v
Line: %v
FuncName: %v

Message:
%v`

	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		file, line := details.FileLine(pc)
		content = fmt.Sprintf(content, level, file, line, details.Name(), msg)
	}

	for _, chatID := range i.recipient {
		msg := tgbotapi.NewMessage(chatID, content)
		i.bot.Send(msg)
	}

}

// Fatal print fatal message to output and quit the application with status 1
func (i *instance) Fatal(v ...interface{}) {
	i.SendLog(FATAL, v)
	os.Exit(1)
}

// Fatalf print formatted fatal message to output and quit the application
// with status 1
func (i *instance) Fatalf(format string, v ...interface{}) {
	i.SendLog(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Error print error message to output
func (i *instance) Error(v ...interface{}) {
	i.SendLog(ERROR, v)
}

// Errorf print formatted error message to output
func (i *instance) Errorf(format string, v ...interface{}) {
	i.SendLog(ERROR, fmt.Sprintf(format, v...))
}

// Warn print warning message to output
func (i *instance) Warn(v ...interface{}) {
	i.SendLog(WARN, v)
}

// Warnf print formatted warning message to output
func (i *instance) Warnf(format string, v ...interface{}) {
	i.SendLog(WARN, fmt.Sprintf(format, v...))
}

// Info print informational message to output
func (i *instance) Info(v ...interface{}) {
	i.SendLog(INFO, v)
}

// Infof print formatted informational message to output
func (i *instance) Infof(format string, v ...interface{}) {
	i.SendLog(INFO, fmt.Sprintf(format, v...))
}

// Debug print debug message to output if debug output enabled
func (i *instance) Debug(v ...interface{}) {
	if i.debug {
		i.SendLog(DEBUG, v)
	}
}

// Debugf print formatted debug message to output if debug output enabled
func (i *instance) Debugf(format string, v ...interface{}) {
	if i.debug {
		i.SendLog(DEBUG, fmt.Sprintf(format, v...))
	}
}
