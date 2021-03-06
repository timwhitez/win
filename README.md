# Win

<a href="https://www.buymeacoffee.com/mjwhitta">🍪 Buy me a cookie</a>

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/mjwhitta/win)](https://goreportcard.com/report/gitlab.com/mjwhitta/win)

## What is this?

This Go module wraps WinHttp and WinInet functions and includes HTTP
clients that use those functions to make HTTP requests. Hopefully,
this makes it easier to make HTTP requests in Go on Windows. WinHttp
can even handle NTLM authentication automatically for you.

Microsoft recommends WinInet over WinHttp unless you're writing a
Windows service. I haven't yet found a way to get NTLM auth working
with WinInet, so I would recommend WinHttp for now.

**Note:** This is probably beta quality at best.

## How to install

Open a terminal and run the following:

```
$ go get -u gitlab.com/mjwhitta/win
```

## Usage

Minimal example:

```
package main

import (
    "io/ioutil"

    "gitlab.com/mjwhitta/log"
    "gitlab.com/mjwhitta/win/winhttp/http" // Chose if you need NTLM
    //"gitlab.com/mjwhitta/win/wininet/http"
)

func main() {
    var b []byte
    var dst = "http://127.0.0.1:8080/asdf"
    var e error
    var headers = map[string]string{
        "User-Agent": "testing, testing, 1, 2, 3...",
    }
    var res *http.Response

	http.DefaultClient.TLSClientConfig.InsecureSkipVerify = true

    if _, e = http.Get(dst, headers); e != nil {
        panic(e)
    }

    if res, e = http.Post(dst, headers, []byte("test")); e != nil {
        panic(e)
    }

    if res.Body != nil {
        if b, e = ioutil.ReadAll(res.Body); e != nil {
            panic(e)
        }
    }

    log.Info(res.Status)
    for k, vs := range res.Header {
        for _, v := range vs {
            log.SubInfof("%s: %s", k, v)
        }
    }
    if len(b) > 0 {
        log.Good(string(b))
    }
}
```

## Links

- [Source](https://gitlab.com/mjwhitta/win)

## TODO

- WinHttp
    - `CONNECT`
    - `DELETE`
    - `OPTIONS`
    - `PATCH`
    - `TRACE`
- WinInet
    - FTP client
    - HTTP client
        - `CONNECT`
        - `DELETE`
        - `OPTIONS`
        - `PATCH`
        - `TRACE`
