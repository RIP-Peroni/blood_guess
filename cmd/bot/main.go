// cmd/bot/main.go
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

	// Создаем GameFinder
	gameFinder := usecases.NewGameFinder(uow.GameRepo)

	createGameUseCase := usecases.NewCreateGameUseCase(uow.GameRepo)
	addPlayersUseCase := usecases.NewAddPlayersUseCase(uow.GameRepo)
	copyPlayersUseCase := usecases.NewCopyPlayersUseCase(uow.GameRepo)
	openPredictionsUseCase := usecases.NewOpenPredictionsUseCase(uow.GameRepo)
	closePredictionsUseCase := usecases.NewClosePredictionsUseCase(uow.GameRepo)
	submitPredictionUseCase := usecases.NewSubmitPredictionUseCase(
		uow.GameRepo,
		uow.UserRepo,
		uow.PredictionRepo,
	)
	startGameUseCase := usecases.NewStartGameUseCase(uow.GameRepo)
	finishGameUseCase := usecases.NewFinishGameUseCase(uow.GameRepo)
	setRealRoleUseCase := usecases.NewSetRealRoleUseCase(uow.GameRepo)
	scoringService := services.NewBasicScoringRules()
	awardPointsUseCase := usecases.NewAwardPointsUseCase(
		uow.GameRepo,
		uow.UserRepo,
		uow.PredictionRepo,
		scoringService,
	)

	botAPI := bot.GetAPI()

	availableCommands := map[string]string{
		"start":       "Начать работу с ботом",
		"help":        "Показать список команд",
		"newgame":     "Создать новую игру (если нет других активных игр)",
		"games":       "Показать активные игры",
		"gameinfo":    "Показать информацию об игре",
		"addplayers":  "Добавить несколько игроков в последнюю созданную игру",
		"copyplayers": "Скопировать игроков из последней завершенной игры",
		"setrealrole": "Установить реальные роли игрокам после игры",
		"openpred":    "Открыть прогнозы для последней созданной игры",
		"closepred":   "Закрыть прогнозы для последней игры с открытыми прогнозами",
		"predict":     "Сделать прогноз на роли игроков в последней игре с открытыми прогнозами",
		"mypredict":   "Показать мои прогнозы",
		"profile":     "Показать мой профиль",
		"startgame":   "Начать реальную игру (после закрытия прогнозов)",
		"finish":      "Завершить последнюю игру в процессе",
		"awardpoints": "Начислить очки за прогнозы (после установки всех ролей)",
	}

	bot.RegisterHandler(handlers.NewStartHandler(botAPI))
	bot.RegisterHandler(handlers.NewHelpHandler(botAPI, availableCommands))
	bot.RegisterHandler(handlers.NewNewGameHandler(botAPI, createGameUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewGamesHandler(botAPI, gameFinder))
	bot.RegisterHandler(handlers.NewGameInfoHandler(botAPI, uow.GameRepo, uow.UserRepo, uow.PredictionRepo))
	bot.RegisterHandler(handlers.NewAddPlayersHandler(botAPI, addPlayersUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewCopyPlayersHandler(botAPI, copyPlayersUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewOpenPredictionsHandler(botAPI, openPredictionsUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewClosePredictionsHandler(botAPI, closePredictionsUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewPredictHandler(botAPI, submitPredictionUseCase, gameFinder, uow.UserRepo))
	bot.RegisterHandler(handlers.NewStartGameHandler(botAPI, startGameUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewFinishGameHandler(botAPI, finishGameUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewSetRealRoleHandler(botAPI, setRealRoleUseCase, gameFinder))
	bot.RegisterHandler(handlers.NewAwardPointsHandler(botAPI, awardPointsUseCase, gameFinder))

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
