// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT

package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var client = &http.Client{}

func (srv ApplicationServer) relay(w http.ResponseWriter, r *http.Request) {
	log.Println("relay", "received request")

	URL, _ := url.Parse(srv.destinationURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = singleJoiningSlash(URL.Path, r.URL.Path)
	r.RequestURI = ""

	r.Header.Set("X-Forwarded-For", r.RemoteAddr)

	log.Println("relay", "sending relayed request to destination", r.URL)
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	log.Println("relay", "copy response body")
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
