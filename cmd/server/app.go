package server

import (
	"log"
	"os"
	"shortUrl/shorten_url/internal/repository"
	"shortUrl/shorten_url/internal/service"
	"shortUrl/shorten_url/pkg/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
)

func Run() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка при загрузке файла .env:", err)
	}

	r := gin.Default()

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	urlRepo := repository.NewInMemoryURLRepository()
	urlService := service.NewUrlService(urlRepo)
	urlHandler := http.NewURLHandler(urlService)

	r.POST("/shorten", urlHandler.ShortenURL)
	r.GET("/shorten/:id", urlHandler.GetURL)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Ошибка при получении обновлений:", err)
	}

	// Обрабатываем каждое обновление
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text != "" {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					// Отправляем сообщение с кнопками
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Сократить ссылку"),
						),
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Сгенерировать QR-код"),
						),
					)
					bot.Send(msg)
				//case "short":
				//	qrCodeFilePath := "qr_code.png"
				//	qrCodeMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, qrCodeFilePath)
				//	bot.Send(qrCodeMsg)
				default:
					// Отвечаем на неизвестные команды
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю такой команды.")
					bot.Send(msg)
				}
			} else {
				// Проверяем, является ли сообщение URL
				if strings.HasPrefix(update.Message.Text, "http://") || strings.HasPrefix(update.Message.Text, "https://") {
					// Создаем сокращенную ссылку
					shortenedURL, err := urlService.Create(update.Message.Text)
					if err != nil {
						log.Printf("Ошибка при сокращении URL: %v", err)
						continue
					}

					// Генерируем QR-код для сокращенной ссылки
					qrCodeFilePath := "qr_code.png"
					err = qrcode.WriteFile(shortenedURL.Shortened, qrcode.Medium, 256, qrCodeFilePath)
					if err != nil {
						log.Printf("Ошибка при генерации QR-кода: %v", err)
						continue
					}

					// Отправляем сокращенную ссылку и QR-код в ответ
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сокращенная ссылка: "+shortenedURL.Shortened)
					bot.Send(msg)

				} else {
					// Отвечаем на сообщение бота
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ссылку для сокращения. Пример: https:// или http://")
					bot.Send(msg)
				}
			}
		}
	}

	return r.Run(":8080")
}
