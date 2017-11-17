# Description

Adds a few columns at the end of each row.

## Usage

```sh
cd auto-agent/util/add_coding_fields
go build add_coding_fields.go
cat  ~/Downloads/userLogs.csv | ./add_coding_fields
```

or pipe it to a file

```sh
cat ~/Downloads/userLogs.csv | ./add_coding_fields > out.csv
```
