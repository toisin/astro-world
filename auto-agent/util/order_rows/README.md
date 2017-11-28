# Description

Orders the rows with the common three rooms/groups first

## Usage

```sh
cd auto-agent/util/order_rows
go build order_rows.go
cat  ~/Downloads/userLogs.csv | ./order_rows
```

or pipe it to a file

```sh
cat ~/Downloads/userLogs.csv | ./order_rows > out.csv
```
