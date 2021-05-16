# Backend security

## Why?

The speed-test server code was assigned to be running locally. However, since the bot and the backend may run separately, it is necessary to apply a security patch to the backend.

The patch will detect the `secret.json` file and switch behavior to ensure code compatibility. You can consider the file as a trigger. If the `secret.json` file does not exist, the webserver will assume that protection is not turn on. And all the functionality is still the same.

## How to apply patch

```shell script
# For security you can apply my patch
cp /PATH/TO/SPT_BOT/patch.diff /PATH/TO/SSRSpeed
cd /PATH/TO/SSRSpeed
# Check if the patch can be applied
git apply --check patch.diff
git apply patch.diff
```

## How to make security work

- Create a `secret.json` file the web project path

```shell script
cd /PATH/TO/SSRSpeed
vim secret.json
```

- Add token or ip in it

## `secret.json` format

See provided [`secret_example.json`](https://github.com/Avimitin/SSRSpeed/blob/master/secret_example.json) .

## Authentication

- For all GET requests, source IP needs to be in the allowed list.
- For all POST requests, source IP or token need to match.
