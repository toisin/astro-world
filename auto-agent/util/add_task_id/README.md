# Description

Reads a CSV file and adds a row called TaskId based on the QuestionText.

## Usage

```sh
cd auto-agent/util/add_task_id
go build
./add_task_id ~/Downloads/userLogs.csv
```

or pipe it to a file

```sh
./add_task_id ~/Downloads/userLogs.csv > out.csv
```
