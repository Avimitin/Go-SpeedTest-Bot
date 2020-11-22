# Backend security

## Why?

To have a `secret.json` file that can be picked up by the app to have access control over who is able to access the web api. If `secret.json` is not found, is considered as unprotected. But all the function will not be affect. So you can choose if you don't want the patch.

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