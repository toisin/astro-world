# Description

Filters a CSV file based on (currently) hardcoded list of PrompIds.

## Usage

```sh
cd auto-agent/util/filter_csv
go build filter_csv.go
./filter_csv ~/Downloads/userLogs.csv
```

or pipe it to a file

```sh
./filter_csv ~/Downloads/userLogs.csv > out.csv
```
