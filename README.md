# Go-SpeedTest-Bot

## Info

This project is used to testing server connection with telegram bot support.

## Todo

- [x] Connect to test backend
- [x] Bot support
- [x] Schedule jobs
- [x] Alert

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

> For more about how I define my security please read [backend_security.md](https://github.com/Avimitin/Go-SpeedTest-Bot/blob/master/docs/backend_security.md)

```shell script
git clone https://github.com/NyanChanMeow/SSRSpeed.git
# For security you can apply my patch
cp /PATH/TO/SPT_BOT/patch.diff /PATH/TO/SSRSpeed
cd /PATH/TO/SSRSpeed
git apply --check patch.diff
git apply patch.diff
python3 web.py
```

## Usage

- Define your settings in `config` directory.

- About bot command please read [bot_command.md](https://github.com/Avimitin/Go-SpeedTest-Bot/blob/master/docs/bot_command.md)
