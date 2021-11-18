# Arraybot: Nautlius

Nautilus handles all command execution for [Arraybot](https://github.com/Arraying/Arraybot).

### Technologies

Server-side:
* Go
* discordgo
* discord-interactions-go

Data:
* PostgreSQL

### Environment variables

These are the environment variables required to run the application properly:
* `AUTH_TOKEN`: the Arraybot token
* `BOT_TOKEN`: the application's bot's token
* `APP_ID`: the application's ID
* `PUBKEY`: the application's public key
* `SERVER`: the development server's ID, blank for production
* `ADMINS`: a semicolon separated list of admin user IDs
