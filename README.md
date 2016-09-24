# onhub

Various utilities for the OnHub router. Tested with the TP-Link variety.

# Installation

```sh
$ go get github.com/benmanns/onhub/cmd/onhubdump
```

# Running

```sh
$ onhubdump
```

This returns a JSON dump of the data from `http://192.168.86.1/api/v1/diagnostic-report`.

If you use a different subnet for your router:

```sh
$ onhubdump http://192.168.85.1/api/v1/diagnostic-report
```

If you want to run against a local, already downloaded dump:

```sh
$ onhubdump path/to/diagnostic-report
```
