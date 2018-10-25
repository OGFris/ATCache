//          ATCache  Copyright (C) 2018  AnimeTwist
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"
)

var Instance Server

type Server struct {
	Http         *http.Server
	ReverseProxy *httputil.ReverseProxy
	ProxyServer  *httptest.Server
}

func (s *Server) Start(port string) {
	s.Http = &http.Server{
		Handler:      &Router{},
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	link, err := url.Parse(URL)
	if err != nil {
		panic(err)
	}

	s.ReverseProxy = httputil.NewSingleHostReverseProxy(link)
	// FIXME:
	// 	Couldn't intercept the response from the reverse proxy
	// 	because the response body could be read once only so
	// 	either I save it and the client won't receive anything
	// 	or the client receive it and I don't intercept anything
	// 	I choose the second choice because I could then just
	// 	download it from another request in the background and
	// 	that's what I did.
	s.ProxyServer = httptest.NewServer(s.ReverseProxy)

	go s.Http.ListenAndServe()
}

func (s *Server) Shutdown() {
	s.Http.Close()
	s.ProxyServer.Close()
}
