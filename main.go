package main

import (
	"os"
    "log"
    "time"
    "fmt"

	"github.com/joho/godotenv"
    tele "gopkg.in/telebot.v4"
)

var (
	listening  = false       // Bot's listening state
	logFile    = "chatlog.txt" // File to store messages
	listenChat *tele.Chat       // The chat the bot is listening to
)

// Function to start listening to messages
func startListening(bot *tele.Bot, chat *tele.Chat) {
	if listening {
		bot.Send(chat, "Already listening to this chat.")
		return
	}

	listening = true
	listenChat = chat

	bot.Send(chat, "Started listening to the conversation.")

	// Open the log file for appending messages
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for listening {
		time.Sleep(1 * time.Second) // Poll interval to simulate active listening
	}
}

// Function to stop listening to messages
func stopListening(bot *tele.Bot, chat *tele.Chat) {
	if !listening {
		bot.Send(chat, "Not currently listening.")
		return
	}

	listening = false
	bot.Send(chat, "Stopped listening to the conversation.")
}

// Function to log incoming messages
func logMessage(bot *tele.Bot, message *tele.Message) {
	if listening && message.Chat.ID == listenChat.ID && message.Text != "" {
		// Open the log file to append messages
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Log the message in the format [timestamp] username: message
		logEntry := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), message.Sender.Username, message.Text)
		if _, err := file.WriteString(logEntry); err != nil {
			log.Println("Error writing to log:", err)
		}
	}
}

func main()  {
    godotenv.Load()
    api_key := os.Getenv("TELEGRAM_BOT_TOKEN")

    // Initialize bot with settings
    p := tele.Settings{
        Token:  api_key,
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    }

    bot, err := tele.NewBot(p)
    if err != nil {
        log.Fatal(err)
        return
    }

    // Handle all messages and log them in the background
    bot.Handle(tele.OnText, func(c tele.Context) error {
        user := c.Sender()
        logMessageToFile(user.Username, c.Text())
        return nil
    })

    // Start the bot in the background
    log.Println("Bot is running...")
    bot.Start()
}

