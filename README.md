# Thonkbot

Thonkbot is an easy way to aggregate chat messages and run queries on them. The bot is written in Go.

## Setup

1. Clone the git repo
1. Copy `config.example.json` to `config.json` and fill with the appropriate keys as needed.
1. Import the MySQL dump into your own database. You will need to manually add yourself to the bot_admins table (userlevel is deprecated and irrelevant).
1. Run `go build` to compile the binary
1. Start the bot with `./thonkbot`!
