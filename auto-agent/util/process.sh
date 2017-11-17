#!/bin/sh

# Usage: ./process.sh file.csv

cd add_task_id
go build add_task_id.go
cd ..

cd add_coding_fields
go build add_coding_fields.go
cd ..

cd filter_csv
go build filter_csv.go
cd ..

cat "$1" | \
  add_task_id/add_task_id | \
  add_coding_fields/add_coding_fields | \
  filter_csv/filter_csv filter_csv/prompt-id-filter.csv
