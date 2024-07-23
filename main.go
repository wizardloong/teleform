package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Client creating
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = viper.GetBool("debug") // Setup debug messages

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Creates chanel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// Income
	for update := range updates {
		if update.Message == nil { // Is it message?
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Only "/start" allowed to start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! What is your name?")
			bot.Send(msg)
		} else {
			// Gathering info
			userName := update.Message.Text

			// Storing to csv
			err := saveToCSV(userName)
			if err != nil {
				log.Println("Storing CSV error:", err)
				continue
			}

			response := fmt.Sprintf("Nice to meet you, %s! Your data stored.", userName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			bot.Send(msg)
		}
	}
}

func initConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
		return err
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
		return err
	}

	return nil
}

// saveToCSV stores data to file
func saveToCSV(name string) error {
	fileName := viper.GetString("filename")

	// Gues file exists
	fileExists := fileExists(fileName)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error while file open: %s\n", err)
		return err
	}
	defer f.Close()

	// Создаем писатель CSV
	writer := csv.NewWriter(f)
	defer writer.Flush()

	// If file just created prints header row
	if !fileExists {
		headers := []string{"Name"}
		err := writer.Write(headers)
		if err != nil {
			log.Println("Error while header row writed to CSV:", err)
			return err
		}
	}

	// Write to CSV
	err = writer.Write([]string{name})
	if err != nil {
		log.Println("Error data writing to CSV:", err)
		return err
	}

	return nil
}

// fileExists checks if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
