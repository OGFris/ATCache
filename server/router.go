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
	"fmt"
	"github.com/AnimeTwist/ATCache/cache"
	"io"
	"net/http"
	"os"
	"strings"
)

type Router struct{}

const URL = "https://twist.moe"

func (_ *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/favicon.ico" {
		return
	}

	response, err := http.Get(URL + path)
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		return
	}
	if err != nil {
		panic(err)
	}
	filePath := cache.CacheDir + strings.NewReplacer("/", "_").Replace(strings.Replace(path, "/", "", 1))
	if _, err := os.Stat(filePath); err == nil {
		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		io.Copy(w, f)
	} else {
		w.Header().Set("Location", Instance.ProxyServer.URL + path)
		w.WriteHeader(http.StatusFound)
		go func() {
			f, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}

			defer f.Close()
			defer response.Body.Close()
			written, err := io.Copy(f, response.Body)
			if err != nil {
				panic(err)
			}

			fmt.Println("Finished downloading: ", path, " Size: ", fmt.Sprint(written/1000000), "MB.")
		}()
	}
}
