# translate server
just a server that implements my [translate package](https://github.com/jhbforlife/golang/tree/main/translate). it caches translations (one week) and supported languages (one day) using an SQLite database, and uses cron to help manage their respective tables.

## Dependencies
- [cron](https://github.com/robfig/cron)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)
- [translate](https://github.com/jhbforlife/golang/tree/main/translate)

## Endpoints

### `GET /languages`
writes back the supported languages as JSON.

### `POST /json`
takes in a JSON body with the `Source` language (optional), language to translate `To`, and the `Original` text to translate, and writes back a [`Translation`](https://github.com/jhbforlife/golang/tree/main/translate/README.md#available-structs).

### `GET /translate`
takes in a query with the `source` language (optional), language to translate `to`, and the `original` text to translate, and writes back a [`Translation`](https://github.com/jhbforlife/golang/tree/main/translate/README.md#available-structs).