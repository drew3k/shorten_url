package server

import (
	"github.com/skip2/go-qrcode"
	"log"
	"os"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"shortUrl/shorten_url/internal/service"
	"shortUrl/shorten_url/pkg/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type Server struct {
	router       *gin.Engine
	bot          *tgbotapi.BotAPI
	shortenedURL *domain.URL
}

func NewServer() *Server {
	return &Server{
		router: gin.Default(),
	}
}

func (s *Server) Initialize() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка при загрузке файла .env:", err)
	}

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

	s.bot = bot
	return nil
}

func (s *Server) handleUpdate(update tgbotapi.Update) {
	qrCodeFilePath := "qr_code.png"

	if update.Message == nil {
		return
	}

	var shortenedURL *domain.URL

	if update.Message.Text != "" {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я бот, который сокращает твою ссылку"+
					" и делает QR-код.")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("🔗Сократить ссылку"),
						tgbotapi.NewKeyboardButton("🤯Сгенерировать QR-код"),
						tgbotapi.NewKeyboardButton("📜Все сразу"),
					),
				)
				s.bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🤷🏻‍Я не знаю такой команды.")
				s.bot.Send(msg)
			}
		} else {
			switch update.Message.Text {
			case "🔗Сократить ссылку":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте свою ссылку")
				s.bot.Send(msg)
			case "🤯Сгенерировать QR-код":
				qrCodeMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, qrCodeFilePath)
				s.bot.Send(qrCodeMsg)
				if err := os.Remove(qrCodeFilePath); err != nil {
					log.Printf("Ошибка при удалении QR-кода %v", err)
				}
			case "📜Все сразу":
				if shortenedURL != nil {
					msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала сократите ссылку.")
					s.bot.Send(msg2)
					return
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте свою ссылку")
				s.bot.Send(msg)
			default:
				if strings.HasPrefix(update.Message.Text, "http://") || strings.HasPrefix(update.Message.Text, "https://") {
					urlRepo := repository.NewInMemoryURLRepository()
					urlService := service.NewUrlService(urlRepo)
					shortenedURL, err := urlService.Create(update.Message.Text)
					if err != nil {
						log.Printf("Ошибка при сокращении URL: %v", err)
						return
					}

					// Генерация QR-кода после успешного сокращения ссылки
					err = qrcode.WriteFile(shortenedURL.Shortened, qrcode.Medium, 256, qrCodeFilePath)
					if err != nil {
						log.Printf("Ошибка при генерации QR-кода: %v", err)
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сокращенная ссылка: "+shortenedURL.Shortened)
					s.bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🤷🏻‍Я не знаю такой команды.")
					s.bot.Send(msg)
				}
			}
		}
	}
}

func (s *Server) SetupRoutes() {
	urlRepo := repository.NewInMemoryURLRepository()
	urlService := service.NewUrlService(urlRepo)
	urlHandler := http.NewURLHandler(urlService)

	s.router.POST("/shorten", urlHandler.ShortenURL)
}

func (s *Server) startTelegramBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Error getting updates:", err)
	}

	for update := range updates {
		s.handleUpdate(update)
	}
}

func (s *Server) Run() error {
	if err := s.Initialize(); err != nil {
		return err
	}
	s.SetupRoutes()
	go s.startTelegramBot()
	return s.router.Run(":8080")
}
