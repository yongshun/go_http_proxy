//
// author xiongyongshun
// project go_http_proxy
// version 1.0
// created 16/8/22 00:56
//
package proxy

import (
    "net/http"
    "fmt"
    "io"
    "github.com/yongshun/go_http_proxy/auth"
)

func copyHeaders(dst, src http.Header) {
    for k := range dst {
        dst.Del(k)
    }
    for k, vs := range src {
        for _, v := range vs {
            dst.Add(k, v)
        }
    }
}

func removeProxyHeaders(request *http.Request) {
    // 根据文档提示, 在客户端进行 HTTP 请求时, RequestURI 需要为空.
    request.RequestURI = ""

    // If no Accept-Encoding header exists, Transport will add the headers it can accept
    // and would wrap the response body with the relevant reader.
    request.Header.Del("Accept-Encoding")
    // curl can add that, see
    // http://homepage.ntlworld.com/jonathan.deboynepollard/FGA/web-proxy-connection-header.html
    request.Header.Del("Proxy-Connection")
    request.Header.Del("Proxy-Authenticate")
    request.Header.Del("Proxy-Authorization")
    // Connection, Authenticate and Authorization are single hop Header:
    // http://www.w3.org/Protocols/rfc2616/rfc2616.txt
    // 14.10 Connection
    //   The Connection general-header field allows the sender to specify
    //   options that are desired for that particular connection and MUST NOT
    //   be communicated by proxies over further connections.
    request.Header.Del("Connection")
}

func handleAuth(writer http.ResponseWriter, clientRequest *http.Request, proxy *ProxyHttpServer) bool {
    if proxy.auth != nil {
        if !proxy.auth.CheckAuth(clientRequest) {
            doResponse(writer, auth.BuildBasicUnauthorized(clientRequest, "proxy_realm"))
            return false
        }
    }
    return true
}

func handleHttpProxy(writer http.ResponseWriter, clientRequest *http.Request, proxy *ProxyHttpServer) {
    // http 代理
    if !clientRequest.URL.IsAbs() {
        // 根据 HTTP 代理的规定, 当客户端连接代理服务器时, 请求行中的 URL 必须是完整的, 即使用绝对路径, 例如
        //  GET http://www.google.com HTTP/1.1
        writer.WriteHeader(500)
        writer.Write([]byte(`NOT HTTP PROXY REQUEST`))
        return
    }

    // 处理 request
    removeProxyHeaders(clientRequest)

    fmt.Printf("%s %s %s\n", clientRequest.Method, clientRequest.URL.String(), clientRequest.Proto)
    for k, v := range clientRequest.Header {
        fmt.Printf("%s: %s\n", k, v)
    }
    fmt.Println("Host:", clientRequest.Host)
    fmt.Println("Path:", clientRequest.URL.Path)
    fmt.Println("RequestURI:", clientRequest.RequestURI)

    var response *http.Response

    response, _ = proxy.transport.RoundTrip(clientRequest)

    doResponse(writer, response)
}

func doResponse(writer http.ResponseWriter, response *http.Response) {
    copyHeaders(writer.Header(), response.Header)
    writer.WriteHeader(response.StatusCode)
    io.Copy(writer, response.Body)
}