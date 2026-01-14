package handlers

import (
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ProfileHandler обрабатывает команду /profile
type ProfileHandler struct {
	*BaseHandler
	userRepo       ports.UserRepository
	predictionRepo ports.PredictionRepository
}

func NewProfileHandler(
	bot BotClient,
	userRepo ports.UserRepository,
	predictionRepo ports.PredictionRepository,
) *ProfileHandler {
	return &ProfileHandler{
		BaseHandler:    NewBaseHandler(bot),
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
	}
}

func (h *ProfileHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "profile")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	telegramID := update.Message.From.ID
	user, err := h.userRepo.FindByTelegramID(telegramID)
	if err != nil {
		// Создаем пользователя, если он не найден
		username := update.Message.From.UserName
		if username == "" {
			username = update.Message.From.FirstName
		}
		user = entities.NewUser(telegramID, username)
		if err := h.userRepo.Save(user); err != nil {
			return h.SendHTML(update.Message.Chat.ID,
				"❌ Ошибка при создании профиля. Попробуйте еще раз.")
		}
	}

	// Получаем все прогнозы пользователя
	allPredictions, err := h.predictionRepo.FindAllByUser(user.ID())
	if err != nil {
		// Если метод не реализован, просто показываем базовую информацию
		return h.showBasicProfile(update.Message.Chat.ID, user)
	}

	// Анализируем прогнозы
	stats := h.analyzePredictions(allPredictions)

	return h.showDetailedProfile(update.Message.Chat.ID, user, stats)
}

func (h *ProfileHandler) showBasicProfile(chatID int64, user *entities.User) error {
	profile := fmt.Sprintf(`<b>👤 Ваш профиль</b>

<b>📛 Имя:</b> %s
<b>🆔 Telegram ID:</b> %d
<b>💰 Баланс:</b> %d очков
<b>📅 Зарегистрирован:</b> %s

<i>Сделайте свой первый прогноз, чтобы увидеть больше статистики!</i>`,
		h.EscapeHTML(user.Username()),
		user.TelegramID(),
		user.Balance(),
		user.CreatedAt().Format("02.01.2006 15:04"))

	return h.SendHTML(chatID, profile)
}

func (h *ProfileHandler) showDetailedProfile(chatID int64, user *entities.User, stats *UserStats) error {
	profile := fmt.Sprintf(`<b>👤 Ваш профиль</b>

<b>📛 Имя:</b> %s
<b>🆔 Telegram ID:</b> %d
<b>💰 Баланс:</b> %d очков
<b>📅 Зарегистрирован:</b> %s

<b>📊 Статистика прогнозов:</b>
• Всего прогнозов: %d
• Правильных: %d
• Неправильных: %d
• Ожидают подсчета: %d

<b>🏆 Лучший результат:</b> +%d очков за одну игру
<b>📉 Худший результат:</b> %d очков за одну игру

<i>Продолжайте в том же духе! 🎯</i>`,
		h.EscapeHTML(user.Username()),
		user.TelegramID(),
		user.Balance(),
		user.CreatedAt().Format("02.01.2006 15:04"),
		stats.TotalPredictions,
		stats.CorrectPredictions,
		stats.IncorrectPredictions,
		stats.PendingPredictions,
		stats.BestScore,
		stats.WorstScore)

	return h.SendHTML(chatID, profile)
}

// UserStats содержит статистику пользователя
type UserStats struct {
	TotalPredictions     int
	CorrectPredictions   int
	IncorrectPredictions int
	PendingPredictions   int
	BestScore            int
	WorstScore           int
}

func (h *ProfileHandler) analyzePredictions(predictions []*entities.Prediction) *UserStats {
	stats := &UserStats{}

	if len(predictions) == 0 {
		return stats
	}

	stats.TotalPredictions = len(predictions)
	stats.BestScore = -1000 // Начальное значение
	stats.WorstScore = 1000

	// Группируем по играм для подсчета очков за игру
	gameScores := make(map[string]int)

	for _, pred := range predictions {
		// Проверяем, начислены ли очки
		if points, awarded := pred.PointsAwarded(); awarded {
			if points > 0 {
				stats.CorrectPredictions++
			} else if points < 0 {
				stats.IncorrectPredictions++
			}

			// Обновляем лучший/худший результат
			gameID := string(pred.GameID())
			gameScores[gameID] += points
		} else {
			stats.PendingPredictions++
		}
	}

	// Находим лучший и худший результат
	for _, score := range gameScores {
		if score > stats.BestScore {
			stats.BestScore = score
		}
		if score < stats.WorstScore {
			stats.WorstScore = score
		}
	}

	// Если нет начисленных очков
	if stats.BestScore == -1000 {
		stats.BestScore = 0
		stats.WorstScore = 0
	}

	return stats
}

func (h *ProfileHandler) Command() string {
	return "profile"
}

func (h *ProfileHandler) Description() string {
	return "Показать мой профиль и статистику"
}
