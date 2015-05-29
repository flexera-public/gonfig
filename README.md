gonfig
======

> A simple tool to help with Go (web) service configs.

gonfig reads a given JSON file containing an arbitrary set of configuration settings and produces
Go code that exposes these settings as variables. The name of the package where the
variables are defined is configurable.

For example, given the following JSON:

```json
{
    "api_endpoint": "https://us-3.rightscale.com",
    "port": 8000,
    "worker": {
        "enabled": true,
        "concurrency": 20
    }
}
```
Running:
```bash
gonfig config.json
```
generates a `config.go` file which includes the following snippet:
```go
var (
        ApiEndpoint string
        Port int64
        Worker *WorkerCfg
)

type WorkerCfg struct {
        Concurrency int64 `json:"concurrency"`
        Enabled     bool  `json:"enabled"`
}
```
The file also contains a `Load` function:
```go
func Load(path string) error
```
As you'd expect that function initializes the variables with the values read from the JSON.
Client code can then consume the configuration file with something like:
```go
package main

func main() {
        err := Load("config.json")
        if err != nil {
                panic("Invalid configuration: " + err.Error())
        }
        // Now use ApiEndpoint, Port and Worker variables
}
```
If you'd rather not pollute the main package with config stuff then `gonfig` can be told to
use a different package and output file:
```bash
gonfig -o config/config.go -p cfg
```
The above Generates the file `config.go` in the `config` directory (creating it if necessary).
The generated code defines the `cfg` package so the client code would look like:
```go
func main() {
        err := cfg.Load("config.json")
        // Check err and use cfg.ApiEndpint, cfg.Port etc.
}
```
### Usage
```
gonfig --help
usage: gonfig [<flags>]

Flags:
  --help               Show help.
  -c, --config=CONFIG  path to JSON configuration, defaults to "config.json"
  -o, --out=OUT        path to output file, defaults to "config.go"
  -p, --package="main"  
                       name of go package containing config code
```
