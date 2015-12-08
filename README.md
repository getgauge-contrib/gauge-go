# gauge-go
[![Gauge Badge](https://cdn.rawgit.com/renjithgr/gauge-js/72f332d11f54e16b74aedb875f702643708156f7/Gauge_Badge_1.svg)](http://getgauge.io)

Go language plugin for Thoughtworks Gauge

[![Build Status](https://travis-ci.org/manuviswam/gauge-go.svg?branch=master)](https://travis-ci.org/manuviswam/gauge-go)

To install plugin
Checkout the repository
run the command
```sh
go build
```
copy gauge-go.exe to <gauge plugin directory>/go/1.0.0/bin/ directory
copy go.json file to <gauge plugin directory>/go/1.0.0/ directory

To initialize a project with gauge-go, in an empty directory run:
```sh
gauge --init go
```

Run specs
```sh
gauge specs/
```

## License

![GNU Public License version 3.0](http://www.gnu.org/graphics/gplv3-127x51.png)
Gauge-GO is released under [GNU Public License version 3.0](http://www.gnu.org/licenses/gpl-3.0.txt)