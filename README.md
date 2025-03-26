# Tabbly
[![Go Reference](https://pkg.go.dev/badge/github.com/electrenator/tabbly.svg)](https://pkg.go.dev/github.com/electrenator/tabbly)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrenator/tabbly)](https://goreportcard.com/report/github.com/electrenator/tabbly)

Ever surprised by how large your tab counter gets within your browser? Due to them still being "needed" sometime or when you accumulate some information for an active project? Well this silly application helps to keep track of that usage and gather personal statistics on which browser window has how many tabs open. Could eventually make some pretty graphs from this data.

Currently this application supports Firefox and it's development version however there are plans to support other browsers too ([#26](https://github.com/electrenator/tabbly/issues/26)). Limitations being that the browsers should save session restore files for this application to read. This meaning that browsers like Tor can never be supported.

## CLI arguments
```
Usage of tabbly:
      --db-location string     Override where the db will be saved. Handy in combination with '--import-legacy'
      --import-legacy string   Legacy file to import into application database. Not recommended to import into already existing database files given it doesn't sort imported entries
      --interval uint16        Time between tab checks in seconds (default 60)
  -v, --verbose                Verbose logging output
      --version                Print the application version then exit
```

## How to run
To run this application one needs to build or install the application. The current version currently has to be manually build using `go build .`, which creates an executable which can be run.

Next to that one can also install this application using Go it's install.
```
go install github.com/electrenator/tabbly@latest
``` 
That makes it available at `~/go/bin/tabbly` to execute. One can execute it by calling that or if you have Go in your [PATH variable](https://go.dev/doc/install) one can directly execute it with `tabbly` in your CLI.

**Note**: installing it that way does *requires* `CGO` to be enabled and the `gcc` compiler to be percent. This is a requirement from the [go-sqlite3](https://github.com/mattn/go-sqlite3?tab=readme-ov-file#installation) dependency.

## Automatic starting
Currently the only way to automatically start this application *on Linux* is to add it to cron after install using adding something along the lines of the following to your cron config (`crontab -e`);
```cron
@reboot ~/go/bin/tabbly
```

## Development
Within th application there is a trigger to allow changing between development and production mode. This only changes the file save locations and adds `-dev` to filenames before there extension.  
Dev mode can be compiled using the `-tabs dev` like the following;
```sh
go run -tags dev -v .
```
or
```sh
go build -tags dev
```
