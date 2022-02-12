package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {

	// 127.0.0.1:2002/xxx -> 127.0.0.1:2003/base/xxx
	rs1 := "http://127.0.0.1:2003/base"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}
	log.Println(url1)
	proxy := httputil.NewSingleHostReverseProxy(url1)
	proxy.ModifyResponse = modifyFunc
	log.Println("Start httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

func modifyFunc(resp *http.Response) error {
	if resp.StatusCode != 200 {
		oldPayload, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		newPayload := []byte("hello " + string(oldPayload))
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
		resp.ContentLength = int64(len(newPayload))
		resp.Header.Set("Content-Length", fmt.Sprint(len(newPayload)))
	}
	return nil
}