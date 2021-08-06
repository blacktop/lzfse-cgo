# lzfse-cgo

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/blacktop/lzfse-cgo/Go)
[![GoDoc](https://godoc.org/github.com/blacktop/lzfse-cgo?status.svg)](https://godoc.org/github.com/blacktop/lzfse-cgo) [![GitHub release (latest by date)](https://img.shields.io/github/v/release/blacktop/lzfse-cgo)](https://github.com/blacktop/lzfse-cgo/releases/latest)
![GitHub](https://img.shields.io/github/license/blacktop/lzfse-cgo?color=blue)

> Go bindings for [lzfse](https://github.com/lzfse/lzfse) compression.

---

## Install

```bash
go get github.com/blacktop/lzfse-cgo
```

## Getting Started

```golang
import (
    "io/ioutil"
    "log"

    "github.com/blacktop/lzfse-cgo"
)

func main() {

    dat, err := ioutil.ReadFile("encoded.file")
    if err != nil {
        log.Fatal(fmt.Errorf("failed to read compressed file: %v", err))
    }

    decompressed = lzfse.DecodeBuffer(dat)

    err = ioutil.WriteFile("decoded.file", decompressed, 0644)
    if err != nil {
        log.Fatal(fmt.Errorf("failed to decompress file: %v", err))
    }
}
```

## Credit

- <https://github.com/zchee/lzfse-cgo>

## License

MIT Copyright (c) 2019-2021 blacktop
