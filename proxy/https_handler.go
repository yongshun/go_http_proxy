//
// author xiongyongshun
// project go_http_proxy
// version 1.0
// created 16/8/22 00:59
//
package proxy

import (
    "net/http"
    "io"
    "net"
    "regexp"
    "fmt"
)

var hasPort = regexp.MustCompile(`:\d+$`)

func copyAndClose(dst, src *net.TCPConn) {
    if _, err := io.Copy(dst, src); err != nil {
    }

    dst.CloseWrite()
    src.CloseRead()
}

func handleHttpsProxy(writer http.ResponseWriter, clientRequest *http.Request, proxy *ProxyHttpServer) {

    // 1. 接管客户端和代理服务器的连接.
    hij, ok := writer.(http.Hijacker)
    if !ok {
        panic("httpserver does not support hijacking")
    }

    proxyClientConnect, _, e := hij.Hijack()
    if e != nil {
        panic("Cannot hijack connection " + e.Error())
    }

    host := clientRequest.URL.Host

    // 2. 检查 host 是否包含端口, 一般来说其端口都是443
    if !hasPort.MatchString(host) {
        host += ":80"
    }

    fmt.Printf("Connect to remote host %s\n", host)

    // 3. 代理服务器通过 tcp 连接到远端服务器.
    targetSiteConnect, err := net.Dial("tcp", host)
    if err != nil {
        fmt.Println("error")
        return
    }

    // 4. 当代理服务器与远端服务器连接上后, 代理服务器会给客户端回复200, 表示 http 隧道建立成功.
    proxyClientConnect.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

    // 5. 当 http 隧道建立成功后, 代理服务器的工作就是将远端服务器和客户端之间的数据相互转发.
    tTargetSiteConnect := targetSiteConnect.(*net.TCPConn)
    tProxyClientConnect := proxyClientConnect.(*net.TCPConn)

    go copyAndClose(tTargetSiteConnect, tProxyClientConnect)
    go copyAndClose(tProxyClientConnect, tTargetSiteConnect)
}