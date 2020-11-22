# Go-SpeedTest-Bot

## Info

This project is used to testing server connection with telegram bot support.

## Todo

- [x] Connect to test backend
- [x] Bot support
- [x] Schedule jobs
- [x] Alert
- [ ] User-defined test result format
- [ ] Integrate the default config setting
- [ ] Bot security
- [ ] Persistence of nodes information

## Deploy

- Build

```shell script
git clone https://github.com/Avimitin/Go-SpeedTest-Bot.git
cd go-speedtest-bot
go build -o bin/
```

- Use [release](https://github.com/Avimitin/Go-SpeedTest-Bot/releases/tag/1.0)

```shell script
./go-speed-test-bot-ver
```

- Backend

> You can also apply the `patch.diff` to protect your backend. For more details about what this patch is please read [backend_security.md](https://github.com/Avimitin/Go-SpeedTest-Bot/blob/master/docs/backend_security.md)

```shell script
git clone https://github.com/NyanChanMeow/SSRSpeed.git
python3 web.py
```

## Usage

- Define your settings in `config` directory.

- About bot command please read [bot_command.md](https://github.com/Avimitin/Go-SpeedTest-Bot/blob/master/docs/bot_command.md)
