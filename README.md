# go-logrus-formatters [![Build Status](https://travis-ci.org/FabienM/go-logrus-formatters.svg?branch=master)](https://travis-ci.org/FabienM/go-logrus-formatters)

This repository contains a set of [logrus] formatters.

## Installation

Since this project supports [semver], preferred way of installation is through [dep] or [gomodules]:

```
dep ensure -add github.com/fabienm/go-logrus-formatters
```

or

```
go get github.com/fabienm/go-logrus-formatters
```

## GELF formatter

The [GELF] formatter supports [1.1 payload specification](http://docs.graylog.org/en/2.4/pages/gelf.html#gelf-payload-specification).

Notable features:

* Logrus levels are converted to syslog levels
* Logrus entries times are converted to UNIX timestamps. 
* Logrus entry fields are prefixed with `_`, excepted `version`, `host`, `short_message`, `full_message`, `timestamp` and `level`, allowing override.
 
### Syslog level mapping

| Logrus | Syslog      |
|--------|-------------|
| Panic  | EMERG (0)   |
| Fatal  | CRIT (2)    |
| Error  | ERR (3)     |
| Warn   | WARNING (4) |
| Info   | INFO (6)    |
| Debug  | DEBUG (7)   |

### Usage

```go
package main

import (
	"os"

	"github.com/fabienm/go-logrus-formatters"
	log "github.com/sirupsen/logrus"
)

func init() {
	hostname, _ := os.Hostname()
	// Log as GELF instead of the default ASCII formatter.
	log.SetFormatter(formatters.NewGelf(hostname))
}

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
	log.WithFields(log.Fields{
		"full_message":  "Backtrace here\n\nmore stuff",
		"user_id":      9001,
		"some_info":    "foo",
		"some_env_var": "bar",
	}).Fatal("A short message that helps you identify what is going on")
}
```

Output:

```json
{"_animal":"walrus","_level_name":"INFORMATIONAL","_size":10,"host":"mylaptop","level":6,"short_message":"A group of walrus emerges from the ocean","timestamp":1522937330.7570872,"version":"1.1"}
{"_some_env_var":"bar","_some_info":"foo","_user_id":9001,"_level_name":"CRITICAL","full_message":"Backtrace here\n\nmore stuff","host":"mylaptop","level":2,"short_message":"A short message that helps you identify what is going on","timestamp":1522937330.7573297,"version":"1.1"}
```

## See also

* [GELF]
* [logrus]
* https://github.com/seatgeek/logrus-gelf-formatter
* https://github.com/xild/go-gelf-formatter

[dep]: https://golang.github.io/dep/
[gomodules]: https://github.com/golang/go/wiki/Modules
[semver]: https://semver.org/
[logrus]: https://github.com/sirupsen/logrus
[GELF]: http://docs.graylog.org/en/2.4/pages/gelf.html
