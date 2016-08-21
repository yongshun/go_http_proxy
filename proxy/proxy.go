//
// author xiongyongshun
// project go_http_proxy
// version 1.0
// created 16/8/17 14:27
//
package proxy

import (
    "net/http"
    "github.com/yongshun/go_http_proxy/auth"
)

type ProxyHttpServer struct {
    transport *http.Transport
    auth      *auth.ProxyAuth
}

func (proxy *ProxyHttpServer) ServeHTTP(writer http.ResponseWriter, clientRequest *http.Request) {
    if !handleAuth(writer, clientRequest, proxy) {
        return
    }

    if clientRequest.Method == "CONNECT" {
        // 当是 CONNECT 方法时, 表示需要建立 http 隧道, 这表明是进行 https 代理.
        handleHttpsProxy(writer, clientRequest, proxy)
    } else {
        handleHttpProxy(writer, clientRequest, proxy)
    }
}

func NewProxyHttpServer() *ProxyHttpServer {
    proxy := ProxyHttpServer{
        transport: &http.Transport{},
        auth: &auth.ProxyAuth{
            User: "xys",
            Password: "xays",
        },
    }
    return &proxy
}