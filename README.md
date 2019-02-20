# HLK-SW16 server
[![Build Status](https://www.travis-ci.org/vitaly-kashtalyan/go_relays_switch.svg?branch=master)](https://www.travis-ci.org/vitaly-kashtalyan/go_relays_switch)

## Description:
You can manage on the hlk-sw16 via a simple server. 

### Install package:
``` bash
go get -u github.com/vitaly-kashtalyan/go_relays_switch
```
You can also manually git clone the repository to:
``` bash
$GOPATH/src/github.com/vitaly-kashtalyan/go_relays_switch
```

### Docker:
Execute the command for build and run container:
```bash
docker build . -t go-relays-switch
docker run -i -t -p 8082:8082 --restart always go-relays-switch
```
or

```bash
$ ./run.sh
```

