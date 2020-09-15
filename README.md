# OctoPrint Telegraf Plugin

You can build this project like so:

```
$ env GOOS=linux GOARCH=arm GOARM=5 go build -o octoprint cmd/main.go
```

To integrate with telegraf, update the telegraf.conf file by adding:

```
[[inputs.execd]]
  command = ["/path/to/octoprintbinary", "-config", "/path/to/plugin.conf"]
  signal = "none"
```

## Helpful resources

* https://github.com/influxdata/telegraf/blob/master/plugins/common/shim/README.md
* https://github.com/influxdata/telegraf/blob/master/EXTERNAL_PLUGINS.md
