# Description

Adds a few columns at the end of each row.

## Usage

```sh
cd auto-agent/util/add_coding_columns
go build add_coding_columns.go
cat  ~/Downloads/userLogs.csv | ./add_coding_columns
```

or pipe it to a file

```sh
cat ~/Downloads/userLogs.csv | ./add_coding_columns > out.csv
```
