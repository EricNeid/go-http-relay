// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT

package server

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (srv ApplicationServer) relay(w http.ResponseWriter, r *http.Request) {
	// Step 1: rewrite URL
	URL, _ := url.Parse(srv.destinationURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = singleJoiningSlash(URL.Path, r.URL.Path)
	r.RequestURI = ""

	// Step 2: adjust Header
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)

	// note: client should be created outside the current handler()
	client := &http.Client{}
	// Step 3: execute request
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: copy payload to response writer
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
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
