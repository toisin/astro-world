# Description

Filters a CSV file based on a list of field names. The field names are provided in a CSV file, with
the name of the field as the header and each row having the possible value. For example:

```csv
PromptId
two_records.same_record.2.1
single_record.hide_performance.1
```

## Usage

```sh
cd auto-agent/util/filter_csv
go build filter_csv.go
cat ~/Downloads/userLogs.csv | ./filter_csv prompt-id-filter.csv
```

or pipe it to a file

```sh
cat ~/Downloads/userLogs.csv | ./filter_csv prompt-id-filter.csv > out.csv
```
