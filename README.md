# Smartlight

Automate electricity usage via web app.

This project was intended for automating lamp usage, but actually it is capable to automate any electronic devices as long as the relay channel used can handle the Watt usage of the device.

## Prerequisite

+ change adaptor according to what microcontroller or mini computer you used in the [code](smartlight-server/cmd/bot/main.go#L31)
+ set BOT_ID on environment variable of your microcontroller or mini computer for further server and bot communication
+ set RELAY_IN on environment variable of your microcontroller or mini computer to indicate which PIN used for connecting to relay
+ adjust the amount of relay input connected with the driver in the [code](smartlight-server/cmd/bot/main.go#L32)
+ set static IP on each bot (dynamic IP make it hard for server to determine the right address to communicate)
+ map the bot ID and static IP into the [code](smartlight-server/cmd/server/bot.go#L13)
+ match bot system timezone to [scheduler timezone](smartlight-server/cmd/bot/main.go#L25)

## Note

+ currently the server still hardcoded to sync the schedule and setting for every 5 minutes
+ currently the bot still hardcoded to send the usages to server for every 5 minutes
+ to compile the bot executable, you need to adjust the target OS and the target architecture.
for example I use Raspberry PI 3 Model B+ so i need to set:
    + GOOS=linux
    + GOARCH=arm
    + GOARM=7