# Description

Reads a CSV file and adds a row called TaskId based on the QuestionText.

## Usage

```sh
cd auto-agent/util/add_task_id
go build
cat ~/Downloads/userLogs.csv | ./add_task_id
```

or pipe it to a file

```sh
cat ~/Downloads/userLogs.csv | ./add_task_id > out.csv
```
