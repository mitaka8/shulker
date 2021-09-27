# Shulker
Shulker is a Minecraft, IRC, Discord compatible chat bot that allows for bridging chat between multiple chat platforms.

Configuration is done with minimal json

## Minecraft
Writing chat to the minecraft server is done with rcon, so you must enable that in your server.properties.  Set your rcon host as "rcon_host" in the config.json and your password in "rcon_pass".

Reading from the log file can be done either over SFTP or by reading a local file.  Local files can be reached with for example `file:///opt/minecraft/server/logs/latest.log` for a minecraft server located at `/opt/minecraft/server`.  These URLs are supposed to be at "log_file" in the config.json.  To avoid having the hostname of the server set as the default source in another instance, set the "name" option in the config.

## IRC
Set the destination in your config as an URL using the irc or ircs scheme.  The bot will figure out if it needs to use TLS based on the scheme of the URL.  To avoid having the hostname of the IRC server sent as the source to a receiving end, set the "name" option in the config.

## Discord
Set up a bot on discord.com and copy its bot token into "bot_token" in the config.json, copy the channel id from the URL bar in your browser and set it as the "channel_id", make a webhook in the integrations tab in the channel settings and set the url as "webhook" in the config.  The guild name will be communicated as the source name to other services.

## Thanks
 - james4k for your Minecraft RCON client
 - hpcloud for your tail implementation
 - thoj for your IRC library
 - andersfylling for your Discord library