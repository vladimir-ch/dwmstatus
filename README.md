# dwmstatus
Status bar for DWM

Prints CPU load averages, used memory, date and time. The amount of used memory should be the same as printed by `free`.

## Installation

```sh
$ go get github.com/vladimir-ch/dwmstatus
```
Then start `dwmstatus` before `dwm`. For example, on Fedora I find it easiest to install `dwm-user` package and modify
the `/usr/bin/dwm-start` script.
