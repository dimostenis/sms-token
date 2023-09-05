# SMS TOKEN READER (MacOS only)

Read SMS token and save it to clipboard.

## 1. Motivation

It saves me ton of time when connecting to various VPNs, filling SSO creds etc.

## 2. Prerequisites

If any is a **NO**, it will not work for you.

1. Have `go` installed (coz I'm not uploading binaries)
1. Be iPhone user
1. Be Mac user
1. SMS with token fits in one of these 2 templates:
    1. ___Token code: {TOKEN}___
    2. ___Use verification code {TOKEN} for Microsoft authentication.___
    - Want more? Send me your SMS and we add it as template.
1. Your terminal must have full disk access
    - Why? To read messages database, which is in `~/Library/Messages`.
    - How to set it: __System Settings__ ➡ __Privacy & Security__ ➡ __Full Disk Access__
    - Test it. This command shall return some SMS(s):
        ```sh
        📟 make test-db
        ```

## 3. Usage

### 3.1 Build it

```sh
📟 make build
```

### 3.1 Test it

In repo root, run:

```sh
📟 make test-ride

# you get a message like:
  :: token '435359' copied to clipboard (16 seconds old)
# or:
  :: token is too old (31 minutes), try again
```

### 3.2 _(Optional)_ Install it system-wide

```sh
# 1. Install
📟 make install
# - creates shell script $REPO/bin/token
# - creates symlink $REPO/bin/token -> /usr/local/bin/token

# 2. Run it anywhere
📟 token

# 3. Uninstall
📟 make uninstall
# - this will delete symlink /usr/local/bin/token
```

## 4. How does it work?

Using `sqlite3` CLI we read messages DB (synced between iOS+MacOS) and then parse message text + timestamp.
Newest message wins.
