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
- [x] Bot security
- [ ] Persistence of nodes information

## Deploy

- Docker

```bash
docker run -d \
	-v /path/to/config:/data \
	avimitin/go-speedtest-bot:latest
```

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

If you want to use systemd:

```shell script
mv bot.service /etc/systemd/system/
mv go-speed-test-bot-ver /usr/local/bin/tgbot
systemctl start tgbot
# If you want to check journal
journalctl -au tgbot -f
```

- Backend

> It's highly recommended you to apply the `patch.diff` 
to protect your backend and get full `RESTful API` support. 
For more details about what this patch is please read 
[backend_security.md](./docs/backend_security.md)

```shell script
git clone https://github.com/NyanChanMeow/SSRSpeed.git
python3 web.py
# Or use wsgi server
pip3 install gunicon
gunicon -w 2 -b 0.0.0.0:10870 -t 0 web:app --log-level critical
```

## Usage

- Define your settings in `config` directory.

- About bot command please read [bot_command.md](./docs/bot_command.md)

