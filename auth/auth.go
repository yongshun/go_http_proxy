//
// author xiongyongshun
// project go_http_proxy
// version 1.0
// created 16/8/17 14:27
//
package auth

import (
    "bytes"
    "encoding/base64"
    "io/ioutil"
    "net/http"
    "strings"
)

var unauthorizedMsg = []byte("407 Proxy Authentication Required")

type ProxyAuth struct {
    User     string
    Password string
}

func (auth *ProxyAuth) CheckAuth(request *http.Request) bool {
    proxyAuthorizationHeader := "Proxy-Authorization"

    authHeader := strings.SplitN(request.Header.Get(proxyAuthorizationHeader), " ", 2)
    request.Header.Del(proxyAuthorizationHeader)
    if len(authHeader) != 2 || authHeader[0] != "Basic" {
        return false
    }

    userPasswordStr, err := base64.StdEncoding.DecodeString(authHeader[1])
    if err != nil {
        return false
    }

    userPassword := strings.SplitN(string(userPasswordStr), ":", 2)
    if len(userPassword) != 2 {
        return false
    }

    return auth.User == userPassword[0] && auth.Password == userPassword[1]
}

func BuildBasicUnauthorized(req *http.Request, realm string) *http.Response {
    return &http.Response{
        StatusCode:    407,
        ProtoMajor:    1,
        ProtoMinor:    1,
        Request:       req,
        Header:        http.Header{"Proxy-Authenticate": []string{"Basic realm=" + realm}},
        Body:          ioutil.NopCloser(bytes.NewBuffer(unauthorizedMsg)),
        ContentLength: int64(len(unauthorizedMsg)),
    }
}