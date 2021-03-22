# Bot commands

## List of commands

|     Command     |                         Function                         |
| :-------------: | :------------------------------------------------------: |
|    `/start`     |                Return start text message                 |
|     `/ping`     | Return whether the backend can be successfully connected |
|    `/status`    |                 Return backend's status                  |
|   `/read_sub`   |   Return nodes information with given subscription URL   |
|    `/result`    |                 Return the latest result                 |
|   `/run_url`    |         Start a test with given subscription URL         |
|  `/list_subs`   | Return list of subscriptions you store in local configs  |
|   `/schedule`   |                    Manage timed tasks                    |

## More details

### CMD `/result`

This command will only return parsed text message like:

```reStructuredText
台湾<-上海01 [O3][0.5]: | loss: 0.00% | localping: 43.08 ms | googleping: 162.70 ms
台湾<-广东01 [O3][0.2]: | loss: 0.00% | localping: 11.49 ms | googleping: 289.30 ms
台湾<-广东04 [O3][1.0]: | loss: 0.00% | localping: 13.26 ms | googleping: 247.44 ms
台湾<-江苏01 [O2][1.0]: | loss: 0.00% | localping: 44.67 ms | googleping: 330.06 ms
```

Enable and configure pastebin options at config.json to make
result can store in https://pastebin.com.

### CMD `/list_subs`

This command will return subscriptions you define in `config.json` .

### CMD `/schedule`

schedule will read the default profile that define in config.json

## Configuration

### Where

Config should store in `~/.config/spt_bot/config.json`
Or you can specific the env `SPT_CFG_PATH` to the path that store
the config.json.

If you are using docker, bot will read config.json in `/data`, so you
can run command like this:

```bash
docker run -d -v /path/to/config:/data spt-bot:0.1
```

Using volumes can easily maintain and backup your configuration.

### How

Example configuration like below, you can mock it at 
[config.json](../config/config.json):

### References

See [references.md](./references.md)

