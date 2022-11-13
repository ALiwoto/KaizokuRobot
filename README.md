# github.com/AnimeKaizoku/KaizokuRobot

Yet another channel poster bot to send messages with buttons and other formattings for Telegram.

## Available commands

* `/send` : The actual useful send command used for sending messages to channels.
 
* `/add` : The add command for adding new chat id to json.
* `/remove` : The remove command for removing chat id from json.
* `/getchats` : The command which is used to get all chats in json.

## Features

* Supports buttons using your norm `[label](buttonurl:urlhere)`.
 
* Mentioning `{lable}{-id}` in last line of message where -id is chat/channel id will pull information about it (Name, -Id)and post it with message.
* Also supports `{label}{*-id}`which is same as above but will post a button too with chatlink of the chat's or channel's id with label as its text.
* You can even mention `*-id` inside any `buttonurl:` and bot will fetch the id's invite link and use that as url.
* Reply to any image while sending the message and bot will send the message as a caption to the image to the channel.
* You can map your own preffered commands by editing the `RHS` of json in `commands.json` file in root of repo
* Maybe more but I don't remember it well.
 
 
## License

This project is licensed under MIT license.


 
##### Feel free to contribute.
