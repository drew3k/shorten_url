package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
	"log"
	"os"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"shortUrl/shorten_url/internal/service"
	"strings"
	"time"
)

type BotService interface {
	Initialize() error
	StartTelegramBot()
}

type BotAPI struct {
	bot              *tgbotapi.BotAPI
	shortenRequested bool
	url              *domain.URL
	urlsList         domain.ShortenedURLList
	generationMsgID  int
}

func NewBotAPI(bot *tgbotapi.BotAPI) *BotAPI {
	return &BotAPI{bot: bot}
}

const qrCodeFilePath = "qr_code.png"

func (b *BotAPI) Initialize() error {
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

	b.bot = bot
	return nil
}

func (b *BotAPI) HandleUpdate(update tgbotapi.Update) {

	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		b.HandleCommand(update)
		return
	}

	switch update.Message.Text {
	case "Все сразу 📌":
		b.AllAtOnce(update, qrCodeFilePath)
	case "Сократить ссылку 🔗":
		b.RequestLink(update)
	case "Сгенерировать QR-код 📲":
		b.GenerateQRCode(update, qrCodeFilePath)
	default:
		b.ProcessLink(update)
	}
}

func (b *BotAPI) HandleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я бот, который сокращает твою ссылку"+
			" и делает QR-код.")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Все сразу 📌"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Сократить ссылку 🔗"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Сгенерировать QR-код 📲"),
			),
		)
		b.bot.Send(msg)
	default:
		b.UnknownCommand(update)
	}
}

func (b *BotAPI) ProcessLink(update tgbotapi.Update) {
	if strings.HasPrefix(update.Message.Text, "http://") || strings.HasPrefix(update.Message.Text, "https://") {
		urlRepo := repository.NewInMemoryURLRepository()
		urlService := service.NewUrlService(urlRepo)
		url, err := urlService.Create(update.Message.Text)
		if err != nil {
			log.Printf("Ошибка при сокращении URL: %v", err)
			return
		}

		if b.shortenRequested {
			b.url = url
		}

		newUrl := domain.URL{
			Original:  url.Original,
			Shortened: url.Shortened,
		}
		b.urlsList.URLs = append(b.urlsList.URLs, newUrl)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сокращенная ссылка: "+url.Shortened)
		b.bot.Send(msg)
		b.shortenRequested = false
	} else {
		b.UnknownCommand(update)
	}
}

func (b *BotAPI) UnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю такой команды 🤷")
	b.bot.Send(msg)
}

func (b *BotAPI) RequestLink(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте свою ссылку")
	b.bot.Send(msg)
	b.shortenRequested = true
}

func (b *BotAPI) GenerateQRCode(update tgbotapi.Update, qrCodeFilePath string) {
	if b.url != nil {
		err := qrcode.WriteFile(b.url.Shortened, qrcode.Medium, 256, qrCodeFilePath)
		if err != nil {
			log.Printf("Ошибка при генерации QR-кода: %v", err)
		} else {
			b.TextGeneration(update)
			qrCodeMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, qrCodeFilePath)
			b.bot.Send(qrCodeMsg)
			delMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, b.generationMsgID)
			b.bot.Send(delMsg)
			if err := os.Remove(qrCodeFilePath); err != nil {
				log.Printf("Ошибка удаления QR-кода %v", err)
			}
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала сократите ссылку.")
		b.bot.Send(msg)
	}
}

func (b *BotAPI) TextGeneration(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Генерация...🎲")
	sentMsg, _ := b.bot.Send(msg)
	b.generationMsgID = sentMsg.MessageID

	time.Sleep(1 * time.Second)
}

func (b *BotAPI) AllAtOnce(update tgbotapi.Update, qrCodeFilePath string) {
	if b.url != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сокращенная ссылка: "+b.url.Shortened)
		b.bot.Send(msg)

		if err := qrcode.WriteFile(b.url.Shortened, qrcode.Medium, 256, qrCodeFilePath); err != nil {
			log.Printf("Ошибка при генерации QR-кода: %v", err)
		} else {
			b.TextGeneration(update)
			qrCodeMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, qrCodeFilePath)
			b.bot.Send(qrCodeMsg)
			delMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, b.generationMsgID)
			b.bot.Send(delMsg)
			if err := os.Remove(qrCodeFilePath); err != nil {
				log.Printf("Ошибка удаления QR-кода %v", err)
			}
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала сократите ссылку.")
		b.bot.Send(msg)
	}
}

func (b *BotAPI) StartTelegramBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Error getting updates:", err)
	}

	for update := range updates {
		b.HandleUpdate(update)
	}
}
