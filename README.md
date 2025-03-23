# Tabbly
<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/Electrenator/Tabbly.svg)](https://pkg.go.dev/github.com/Electrenator/Tabbly) -->
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/Electrenator/Tabbly)](https://goreportcard.com/report/github.com/Electrenator/Tabbly) -->
<!-- Above not yet available due to not being published -->

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
To run this application one needs to *either* build the *Go version* or run the *Python version*. The Go version currently has to be manually build using `go build ./src-go/main.go`, which creates an executable which can be run. With the Python version requiring some more work.

### Python version
The python version need the following things installed;
- Python 3.8+
- Python 3.8+ PIP

With those installed it's recommended to create an virtual environment and installing the requirements of this application within.

```bash
python -m venv venv
source venv/bin/activate # ← might be different depending on your OS
pip install -r requirements.txt

# Use "deactivate" to exit this virtual environment
```

With that installed the application can be run using `python src/main.py`. The only thing you need to make sure of now is that the `Display current activity as a status message.` setting is enabled within Discord if you want to show tab usages within there. *Do note that the Discord integration feature won't be transferred over to the Go version, which will eventually replace this Python version*.

![Tabbly in use within Discord. The application is displaying the usage of 83 active browser tabs.](https://user-images.githubusercontent.com/18311389/151074155-78ccf239-5127-4e7a-8380-f7038ade6338.png)

Lastly the Python version also has development dependances. THese are only required while developing and give tools like linters and formatters.

```bash
pip install -r requirements.dev.txt
```

## Development
Within the Go version there is a trigger to allow changing between development and production mode. Dev mode can be compiled using the `-tabs dev` like the following;
```sh
go run -tags dev -v ./src-go
```
or
```sh
go build -tags dev ./src-go/main.go
```
