# Description

Adds a few ? to the cells for the new columns we added using add_coding_columns.

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
