# go-lzfse

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/blacktop/go-lzfse/Go)
[![GoDoc](https://godoc.org/github.com/blacktop/go-lzfse?status.svg)](https://godoc.org/github.com/blacktop/go-lzfse) [![GitHub release (latest by date)](https://img.shields.io/github/v/release/blacktop/go-lzfse)](https://github.com/blacktop/go-lzfse/releases/latest)
![GitHub](https://img.shields.io/github/license/blacktop/go-lzfse?color=blue)

> Go bindings for [lzfse](https://github.com/lzfse/lzfse) compression.

---

## Install

```bash
go get github.com/blacktop/go-lzfse
```

## Getting Started

```golang
import (
    "io/ioutil"
    "log"

    "github.com/blacktop/go-lzfse"
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

- <https://github.com/zchee/go-lzfse>

## License

MIT Copyright (c) 2019-2021 blacktop
