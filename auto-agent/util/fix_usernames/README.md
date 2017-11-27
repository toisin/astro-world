# Description

Fixes the usernames so that they all follow the right format

## Usage

```sh
cd auto-agent/util/fix_usernames
go build
cat ~/Downloads/userLogs.csv | ./fix_usernames
```

or pipe it to a file

```sh
cat ~/Downloads/userLogs.csv | ./fix_usernames > out.csv
```
