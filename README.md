# Arraybot: Nautilus

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
* `COMMANDS_SECRET`: The secret token to use to trigger actions from the panel.
* `DEV_SERVER`: The development server's ID, can be blank or missing in production
* `HOST_MONITOR`: The monitor API's hostname
* `PORT_COMMANDS`: The command handler's port
* `PORT_MONITOR`: The monitor API's port
* `SCHEME_MONITOR`: The monitor API's scheme

### Command line flags

These are the flags that can be passed into the command line.

#### --prod

Default value: `false`.
This runs the application in production mode.
As such, no WebSocket connection will be established.

#### --register

Default value: `false`.
This ensures to re-register all slashcommands in Discord.
If `--prod` is set to `true`, this will register them as global commands.
Otherwise, it will register them as local commands in the server in the `DEV_SERVER` envrionment variable.

#### --noserve

Default value: `false`.
If this is set to `true`, this will soley register the commands.
The database will not be connected to, neither will a WebSocket, and the HTTP server will not bind.