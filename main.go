//
// author xiongyongshun
// project go_http_proxy
// version 1.0
// created 16/8/17 14:27
//
package main

import (
    "github.com/yongshun/go_http_proxy/proxy"
    "log"
    "net/http"
)


func main() {
    proxy := proxy.NewProxyHttpServer()
    log.Fatal(http.ListenAndServe(":8080", proxy))
}