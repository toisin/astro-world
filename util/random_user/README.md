# Description

Splits a CSV file into two based on random user (random is using a fixed seed so that it always
generates the same random usernames).

## Usage

```sh
cd auto-agent/util/random_user
go build
cat ~/Downloads/userLogs.csv | ./random_user group1.csv group2.csv
```
