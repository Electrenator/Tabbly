# Tabbly
[![Go Reference](https://pkg.go.dev/badge/github.com/Electrenator/tabbly.svg)](https://pkg.go.dev/github.com/Electrenator/tabbly)
[![Go Report Card](https://goreportcard.com/badge/github.com/Electrenator/tabbly)](https://goreportcard.com/report/github.com/Electrenator/tabbly)

Ever surprised by how large your tab counter gets within your browser? Due to them still being "needed" sometime or when you accumulate some information for an active project? Well this application is to help one keep track of that usage and gather data on which browser window has how many tabs open. Hope to one day create of this.

Currently this application supports Firefox and it's development version. There are plans to support other browsers too (#26) but there currently isn't a timeline on that. Limitations being that the browsers should save session restore files for this application to read.

## CLI arguments
```
Usage of tabbly:
      --db-location string     Override where the db will be saved. Handy in combination with '--import-legacy'
      --dryrun                 Disable file writing
      --import-legacy string   Legacy file to import into application database. Not recommended to import into already existing database files given it doesn't sort imported entries
      --interval uint16        Time between tab checks in seconds (default 60)
  -v, --verbose                Verbose logging output
```

## How to run
To run this application one needs to build or install the application. The current version currently has to be manually build using `go build .`, which creates an executable which can be run.

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
