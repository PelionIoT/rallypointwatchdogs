# RallyPoint Maestro Watchdog Plugins

To use, build the watchdogs using the `./build.sh` script. This will build a `.so` file which is a Go version of a shared libary (aka. "go plugin")

#### in Maestro

Usage in [maestro](https://github.com/armPelionEdge/maestro)'s config file:

```
watchdog:
  path: "/home/ed/work/gostuff/src/github.com/armPelionEdge/rallypointwatchdogs/rp100/rp100wd.so"
  opt1: "/tmp/devOSkeepalive"
  opt2: "30"
```
