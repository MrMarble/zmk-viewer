# ZMK Layout viewer

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mrmarble/zmk-layout-viewer)
[![Vuln](https://github.com/MrMarble/zmk-layout-viewer/actions/workflows/vuln.yml/badge.svg)](https://github.com/MrMarble/zmk-layout-viewer/actions/workflows/vuln.yml)
[![golangci-lint](https://github.com/MrMarble/zmk-layout-viewer/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/MrMarble/zmk-layout-viewer/actions/workflows/golangci-lint.yml)

A **work in progress** Cli tool to generate preview images from a zmk .keymap file.

## Usage

```shell
Usage: zmk-layout-viewer <keyboard>

Arguments:
  <keyboard>    Keyboard name to fetch layout.

Flags:
  -h, --help           Show context-sensitive help.
  -f, --file=STRING    ZMK .keymap file
  -t, --transparent    Use a transparent background.
  -o, --output="."     Output directory.
      --debug          Enable debug logging.
```

Keyboard name should be the same as in https://config.qmk.fm.

```shell
zmk-layout-viewer cradio
```
Will output this image:

![](assets/layout.png)

You can pass a .keymap file ([this one for reference](https://github.com/MrMarble/zmk-config/blob/master/config/cradio.keymap)) to generate the layout with bindings

```shell
zmk-layout-viewer -f ~/zmk-config/config/cradio.keymap cradio
```
will output an image for each layer:

![](assets/default_layer.png)
![](assets/left_layer.png)
![](assets/right_layer.png)
![](assets/tri_layer.png)
