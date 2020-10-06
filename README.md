# OctoPrint Telegraf Plugin

![octoprint](octoprint.png "icon")

Gather useful information from the octoprint API

## Configuration

Example plugin.conf required for plugin

```toml
[[inputs.octoprint]]
url="http://xxx.xxx.x.xxx:xxx/"
apikey=""
```

If you have the [Filament Manager Plugin](https://plugins.octoprint.org/plugins/filamentmanager/) then you can configure the plugin to use the external filament manager database,
[Follow this guide for setup on raspberry pi](https://github.com/malnvenshorn/OctoPrint-FilamentManager/wiki/Setup-PostgreSQL-on-Raspbian-(Stretch)).

Example of an updated plugin.conf to support the postgres database

```toml
[[inputs.octoprint]]
url="http://xxx.xxx.x.xxx:xxx/"
apikey=""
dbnamepsql="octoprint_filamentmanager"
userpsql="octoprint"
passpsql="xxxx"
ip="xxx.xxx.x.xxx"
```

To integrate with telegraf, extend the telegraf.conf using the following example

```toml
[[inputs.execd]]
command = ["/path/to/octoprintbinary", "-config", "/path/to/plugin.conf"]
signal = "none"
```

## Development

Refer to [deploy.sh](deploy.sh) for a building and deploying example

Useful for debugging

```bash
journalctl -l -u telegraf.service -b -n 10
```

## Helpful resources

* [External plugin overview](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim/README.md)
* [Examples of other external plugins](https://github.com/influxdata/telegraf/blob/master/EXTERNAL_PLUGINS.md)
* [Example setting up external Postgres DB For filament manager on Raspberry PI](https://github.com/malnvenshorn/OctoPrint-FilamentManager/wiki/Setup-PostgreSQL-on-Raspbian-(Stretch))

## Example Dashboard

![example](Example.PNG "example")