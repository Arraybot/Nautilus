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
* `ADMINS`: A semicolon separated list of admin user IDs
* `APP_ID`: The application's ID
* `APP_PUBKEY`: The application's public key
* `BOT_TOKEN`: The application's bot's token
* `DEV_SERVER`: The development server's ID, can be blank or missing in production
