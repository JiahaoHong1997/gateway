package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	proxy_addr = "http://127.0.0.1:2003"
	port = "2002"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// step 1 解析代理地址，并更改请求体的协议和主机
	proxy, _ := url.Parse(proxy_addr)
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host
	log.Println(r.URL)

	// step 2 请求下游
	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	// step 3 把下游请求返回给上游
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	bufio.NewReader(resp.Body).WriteTo(w)

}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port " + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}