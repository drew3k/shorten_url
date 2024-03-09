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
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я бот, который сокращает твою ссылку"+
						"и делает QR-код.")
					msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте свою ссылку")
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Сократить ссылку"),
						),
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Сгенерировать QR-код"),
						),
					)
					bot.Send(msg)
					bot.Send(msg2)

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

					// Отправляем сокращенную ссылку
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сокращенная ссылка: "+shortenedURL.Shortened)
					bot.Send(msg)

					msg3 := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите дальнейшее действие:")
					bot.Send(msg3)

				} else {
					// Генерируем QR-код для сокращенной ссылки
					shortenedURL, err := urlService.Create(update.Message.Text)
					if err != nil {
						log.Printf("Ошибка при сокращении URL: %v", err)
						continue
					}
					qrCodeFilePath := "qr_code.png"
					err = qrcode.WriteFile(shortenedURL.Shortened, qrcode.Medium, 256, qrCodeFilePath)
					if err != nil {
						log.Printf("Ошибка при генерации QR-кода: %v", err)
						continue
					}
					// Отправляем QR-код в ответ
					qrCodeMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, qrCodeFilePath)
					bot.Send(qrCodeMsg)
				}
			}
		}
	}

	return r.Run(":8080")
}
