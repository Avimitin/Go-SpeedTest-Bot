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
| `/set_default`  |   Set the default subscriptions URL with given remarks   |
|   `/set_chat`   |       Set the default chat room for sending result       |
| `/set_def_mode` |           Set the default test method and mode           |
|   `/run_def`    |              Run a test with default config              |
|   `/schedule`   |                    Manage timed tasks                    |
| `/set_interval` |         Set interval between each schedule test          |
|   `/set_exin`   |         Set default exclude and include remarks          |
|   `/show_def`   |         Show all the default settings                    |
|   `/add_admin`  |         Grant bot permission to someone                  |

## More details

### CMD `/result`

This command will only return parsed text message like:

```reStructuredText
台湾<-上海01 [O3][0.5]: | ls: 0.00% | lp: 43.08 ms | gp: 162.70 ms
台湾<-广东01 [O3][0.2][FAKEIEPL]: | ls: 0.00% | lp: 11.49 ms | gp: 289.30 ms
台湾<-广东04 [O3][1.0]: | ls: 0.00% | lp: 13.26 ms | gp: 247.44 ms
台湾<-江苏01 [O2][1.0]: | ls: 0.00% | lp: 44.67 ms | gp: 330.06 ms
```

In this case, ls aka loss, lp aka local ping, gp aka google ping.

### CMD `/list_subs`

This command will return subscriptions you define in `configs/subs.ini` .

### CMD `/set_default`

You can only set default URL that define in `configs/subs.ini` .

### CMD `/set_def_mod`

Method usable: `ST_ASYNC | SOCKET | SPEED_TEST_NET | FAST`

Mode usable: `TCP_PING | WEB_PAGE_SIMULATION | ALL`

For more detail, check out backend [README](https://github.com/NyanChanMeow/SSRSpeed#test-modes) .

### CMD `/schedule`

Before starting schedule jobs, you should set default remarks with command `/set_default` and default chat room with command `/set_chat` .

### CMD `/set_chat` 

Command `/set_chat` now only support group or user's id.

### CMD `/add_admin`

Command `/add_admin` can only add user when you reply to someone. Will only add user who have username. 