```mermaid
stateDiagram-v2
    [*] --> CREATED : /newgame
    CREATED --> PREDICTIONS_OPEN : /openpredictions
    PREDICTIONS_OPEN --> PREDICTIONS_CLOSED : /closepredictions или таймер
    PREDICTIONS_CLOSED --> IN_PROGRESS : Начало реальной игры
    IN_PROGRESS --> FINISHED : /finishgame
    FINISHED --> [*] : Завершение
    
    state PREDICTIONS_OPEN {
        [*] --> AWAITING_PREDICTION
        AWAITING_PREDICTION --> PREDICTION_SUBMITTED : Прогноз отправлен
        PREDICTION_SUBMITTED --> AWAITING_PREDICTION : /predict повторно
    }
    
    state FINISHED {
        [*] --> CALCULATING_POINTS
        CALCULATING_POINTS --> SENDING_RESULTS
        SENDING_RESULTS --> [*]
    }
```