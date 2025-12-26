# Bot Testing Instructions

## 1. Obtaining a Bot Token

1. Open Telegram and search for @BotFather
2. Send `/newbot`
3. Follow the instructions:
- Enter the bot name (e.g.: "Blood Predictions Bot")
- Enter the bot username (e.g.: `blood_predictions_bot`)
4. Save the received token

## 2. Project Setup

```bash
# Clone the repository (if you haven't already)
git clone <your-repository>
cd blood_guess

# Install dependencies
make setup

# Edit the .env file
nano .env # or any text editor
# Add the token: TELEGRAM_BOT_TOKEN=your_token