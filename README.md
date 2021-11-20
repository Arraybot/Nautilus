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
* `DEV_SERVER`: the development server's ID, can be blank in production
* `ADMINS`: a semicolon separated list of admin user IDs
