package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"shortUrl/shorten_url/internal/repository"
	"shortUrl/shorten_url/internal/service"
	"shortUrl/shorten_url/pkg/http"
)

func Run() error {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("Failed to load .env file: %v", err)
	//}

	//db, err := pgx.Connect(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
	//	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	//	os.Getenv("DB_NAME"), os.Getenv("SSL_MODE")))
	//if err != nil {
	//	return fmt.Errorf("failed to connect to database: %v", err)
	//}
	//defer db.Close(context.Background())
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

	return r.Run(":8080")

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

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот, который отвечает на все твои сообщения.")

		// Отправляем ответное сообщение
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
		}
	}

	return r.Run(":8080")
}
