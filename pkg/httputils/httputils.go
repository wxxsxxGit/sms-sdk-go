package httputils

import (
	"fmt"
	"net/http"
	"strings"
)

func CurlStyleOutput(req *http.Request, resp *http.Response, b1 []byte, b2 []byte) string {
	method := req.Method
	proto := req.Proto
	url := req.URL
	// host := req.Host
	// fmt.Printf("%s\n",url.)
	sb := &strings.Builder{}
	sb.WriteString(strings.Repeat("=", 20) + "\n")
	sb.WriteString("http请求\n")
	sb.WriteString(fmt.Sprintf("%s %s %s\n", method, url.Path, proto))
	sb.WriteString(fmt.Sprintf("Host: %s\n", req.Host))
	for k, v := range req.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	sb.WriteString("\n")
	sb.WriteString(string(b1))
	sb.WriteString("\n")

	respProto := resp.Proto
	respStatus := resp.Status
	sb.WriteString("http响应\n")
	sb.WriteString(fmt.Sprintf("%s %s\n", respProto, respStatus))
	for k, v := range resp.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	sb.WriteString("\n")
	sb.WriteString(string(b2))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 20) + "\n")
	return sb.String()
}
