package main

import (
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/domain/services"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := telegram.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	bot, err := telegram.NewBot(config)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	uow := persistence.NewUnitOfWork()

	createGameUseCase := usecases.NewCreateGameUseCase(uow.GameRepo)
	openPredictionsUseCase := usecases.NewOpenPredictionsUseCase(uow.GameRepo)
	closePredictionsUseCase := usecases.NewClosePredictionsUseCase(uow.GameRepo)
	submitPredictionUseCase := usecases.NewSubmitPredictionUseCase(
		uow.GameRepo,
		uow.UserRepo,
		uow.PredictionRepo,
	)
	addPlayerUseCase := usecases.NewAddPlayerUseCase(uow.GameRepo)
	startGameUseCase := usecases.NewStartGameUseCase(uow.GameRepo)
	scoringService := services.NewBasicScoringRules()
	finishGameUseCase := usecases.NewFinishGameUseCase(
		uow.GameRepo,
		uow.UserRepo,
		uow.PredictionRepo,
		scoringService,
	)

	botAPI := bot.GetAPI()

	availableCommands := map[string]string{
		"start":     "Начать работу с ботом",
		"help":      "Показать список команд",
		"newgame":   "Создать новую игру",
		"games":     "Показать активные игры",
		"gameinfo":  "Показать информацию об игре",
		"addplayer": "Добавить игрока в игру",
		"openpred":  "Открыть прогнозы для игры",
		"closepred": "Закрыть прогнозы для игры",
		"predict":   "Сделать прогноз на игру",
		"mypredict": "Показать мои прогнозы",
		"profile":   "Показать мой профиль",
		"startgame": "Начать реальную игру (после закрытия прогнозов)",
		"setrole":   "Установить реальную роль игрока",
		"finish":    "Завершить игру и подсчитать результаты",
	}

	bot.RegisterHandler(handlers.NewStartHandler(botAPI))
	bot.RegisterHandler(handlers.NewHelpHandler(botAPI, availableCommands))
	bot.RegisterHandler(handlers.NewNewGameHandler(botAPI, createGameUseCase))
	bot.RegisterHandler(handlers.NewGamesHandler(botAPI, uow.GameRepo))
	bot.RegisterHandler(handlers.NewGameInfoHandler(botAPI, uow.GameRepo, uow.UserRepo, uow.PredictionRepo))
	bot.RegisterHandler(handlers.NewAddPlayerHandler(botAPI, addPlayerUseCase, uow.GameRepo))
	bot.RegisterHandler(handlers.NewOpenPredictionsHandler(botAPI, openPredictionsUseCase, uow.GameRepo))
	bot.RegisterHandler(handlers.NewClosePredictionsHandler(botAPI, closePredictionsUseCase, uow.GameRepo))
	bot.RegisterHandler(handlers.NewPredictHandler(botAPI, submitPredictionUseCase, uow.GameRepo, uow.UserRepo))
	bot.RegisterHandler(handlers.NewStartGameHandler(botAPI, startGameUseCase, uow.GameRepo))
	bot.RegisterHandler(handlers.NewFinishGameHandler(botAPI, finishGameUseCase, uow.GameRepo))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting bot...")
		if err := bot.Start(); err != nil {
			log.Printf("Bot stopped with error: %v", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	bot.Stop()
	log.Println("Application stopped")
}
