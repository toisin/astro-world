#!/bin/sh

if [ "$#" -ne 5 ]; then
  echo "Usage: $0 source.csv rm2g1.csv rm2g2.csv rm10g1.csv rm10g2.csv" >&2
  exit 1
fi

cd add_task_id
go build add_task_id.go
cd ..

cd add_coding_columns
go build add_coding_columns.go
cd ..

cd add_coding_fields
go build add_coding_fields.go
cd ..

cd filter_csv
go build filter_csv.go
cd ..

cd order_rows
go build order_rows.go
cd ..

cd random_user
go build random_user.go
cd ..

cat "$1" | \
  add_task_id/add_task_id | \
  add_coding_columns/add_coding_columns | \
  add_coding_fields/add_coding_fields | \
  filter_csv/filter_csv filter_csv/prompt-id-filter.csv | \
  order_rows/order_rows | \
  random_user/random_user "$2" "$3" "$4" "$5"
