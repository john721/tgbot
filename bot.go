package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func separateCommand(input string) (string, string) {
	index := strings.Index(input, " ")
	if index == -1 {
		return input, ""
	}
	command := input[:index]
	rest := input[index+1:]

	return command, rest
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Create a buffered writer for the file

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		fmt.Printf("User: %v\n", update.FromChat().FirstName)
		filename := update.FromChat().FirstName + "_food.txt"
		foods := []string{}
		file, err := os.ReadFile(filename)
		if err == nil {
			lines := strings.Split(string(file), "\n")
			// Remove empty lines
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					foods = append(foods, line)
				}
			}
		}
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}
		responseText := ""
		cmd, parameters := separateCommand(update.Message.Text)
		switch strings.ToLower(cmd) {
		case "clear":
			foods = []string{}
			responseText = "OK, 已刪除清單"
		case "list":
			if len(foods) == 0 {
				responseText = "--食物清單是空的--\n"
				break
			}
			responseText = "--食物清單--\n"
			for i, food := range foods {
				responseText += strconv.Itoa(i) + ". " + food
				responseText += "\n"
			}
		case "food":
			foods = append([]string{parameters}, foods...)
			responseText = "OK, 已加入 " + parameters
			fmt.Printf("foods=%v\n", foods)
			// Write each string in the slice to the file
		case "choose":
			rand.Seed(time.Now().UnixNano())
			randomNum := 0
			fmt.Printf("foods=%v\n", foods)
			if len(foods) == 0 {
				responseText = "你還沒加過食物哦!"
				break
			} else if len(foods) <= 2 {
				// Generate a random number between 0 and len(food) (exclusive)
				randomNum = rand.Intn(len(foods))
			} else {
				// Generate a random number between 2 and len(food) (exclusive)
				randomNum = rand.Intn(len(foods)-2) + 2
			}
			fmt.Printf("randonNum=%v\n", randomNum)
			foodToday := foods[randomNum]
			if randomNum == len(foods) {
				newfoods := []string{foodToday}
				foods = append(newfoods, foods[:randomNum]...)
			} else {
				newfoods := []string{foodToday}
				newfoods = append(newfoods, foods[:randomNum]...)
				foods = append(newfoods, foods[randomNum+1:]...)
			}

			responseText = "就決定是你了 -- " + foodToday + "!"
		default:
			responseText = "usage: \n1. food FOOD_NAME\n2. choose\n3. list\n4. clear"
		}

		// write all foods to file
		os.WriteFile(filename, []byte(strings.Join(foods, "\n")), 0777)

		// Now that we know we've gotten a new message, we can construct a
		// reply! We'll take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		msg.ReplyToMessageID = update.Message.MessageID

		replyKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Choose"),
				tgbotapi.NewKeyboardButton("List"),
				// tgbotapi.NewKeyboardButton("Clear"),
				tgbotapi.NewKeyboardButton("Help"),
			),
		)
		replyKeyboard.OneTimeKeyboard = true
		msg.ReplyMarkup = replyKeyboard

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}
}
