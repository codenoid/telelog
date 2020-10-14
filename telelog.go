package telelog

import (
	"fmt"
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

// Instance telelog instance
type Instance struct {
	name      string
	debug     bool
	token     string
	file      string
	recipient []int64

	osLogger *log.Logger

	bot *tgbotapi.BotAPI
}

// LoggerNew create new telelog Instance
func LoggerNew() *Instance {

	i := &Instance{}

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
		i.SetRecipient(path)
	}

	return i
}

func (i *Instance) SetToken(token string) {
	i.token = token
}

func (i *Instance) SetAppName(name string) {
	i.name = name
}

func (i *Instance) SetDebug(debug bool) {
	i.debug = debug
}

func (i *Instance) SetRecipient(files ...string) {
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

func (i *Instance) Connect() error {
	bot, err := tgbotapi.NewBotAPI(i.token)
	if err != nil {
		return err
	}

	i.bot = bot

	return nil
}

func (i *Instance) sendLog(level string, msg string) {

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

	for _, chatID := range i.recipient {
		msg := tgbotapi.NewMessage(chatID, content)
		i.bot.Send(msg)
	}

}

// Fatal print fatal message to output and quit the application with status 1
func (i *Instance) Fatal(v ...interface{}) {
	i.sendLog(FATAL, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf print formatted fatal message to output and quit the application
// with status 1
func (i *Instance) Fatalf(format string, v ...interface{}) {
	i.sendLog(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Error print error message to output
func (i *Instance) Error(v ...interface{}) {
	i.sendLog(ERROR, fmt.Sprintln(v...))
}

// Errorf print formatted error message to output
func (i *Instance) Errorf(format string, v ...interface{}) {
	i.sendLog(ERROR, fmt.Sprintf(format, v...))
}

// Warn print warning message to output
func (i *Instance) Warn(v ...interface{}) {
	i.sendLog(WARN, fmt.Sprintln(v...))
}

// Warnf print formatted warning message to output
func (i *Instance) Warnf(format string, v ...interface{}) {
	i.sendLog(WARN, fmt.Sprintf(format, v...))
}

// Info print informational message to output
func (i *Instance) Info(v ...interface{}) {
	i.sendLog(INFO, fmt.Sprintln(v...))
}

// Infof print formatted informational message to output
func (i *Instance) Infof(format string, v ...interface{}) {
	i.sendLog(INFO, fmt.Sprintf(format, v...))
}

// Debug print debug message to output if debug output enabled
func (i *Instance) Debug(v ...interface{}) {
	if i.debug {
		i.sendLog(DEBUG, fmt.Sprintln(v...))
	}
}

// Debugf print formatted debug message to output if debug output enabled
func (i *Instance) Debugf(format string, v ...interface{}) {
	if i.debug {
		i.sendLog(DEBUG, fmt.Sprintf(format, v...))
	}
}
