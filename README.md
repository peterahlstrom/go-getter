# go-getter

## Prerequisites
Go version > 1.23.5

## Usage

In root directory:

Build application binary
```bash
$ go build main.go
```

Specify endpoints in `config.json`

Run the application
```bash
$ ./main <port>
```

### config.json example

```json
{
  "logPath": "go-getter.log",
  "concurrentScriptsLimit": 3,
  "endpoints": {
    "/data1": {
      "scriptPath": "./script1.sh",
      "requireAuth": true,
      "apiKeys": {
        "abc123": "dev",
        "def456": "prod"
      }
    },
    "/data2": {
      "scriptPath": "./script2.sh",
      "requireAuth": false
    }
  }
}
```
