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
	"github.com/AnimeTwist/ATCache/cache"
	"github.com/AnimeTwist/ATCache/cache/queue"
	"github.com/OGFris/Treagler"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Router struct{}

var URL string

// TODO: Make separated functions for each process to make it more clear.
func (_ *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fail := func() error {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	if path == "/favicon.ico" {
		if _, err := os.Stat(cache.Dir + "favicon.ico"); err != nil {
			response, err := http.Get(URL + path)
			Treagler.Do(func() error { return err }, fail, URL+path)

			defer response.Body.Close()
			bytes, err := ioutil.ReadAll(response.Body)
			Treagler.Do(func() error { return err }, fail, URL+path)

			f, err := os.Create(cache.Dir + "favicon.ico")
			Treagler.Do(func() error { return err }, fail, URL+path)
			f.Write(bytes)
			w.Write(bytes)
		} else {
			f, err := os.Open(cache.Dir + "favicon.ico")
			Treagler.Do(func() error { return err }, fail, URL+path)

			io.Copy(w, f)
		}
	} else {

		paths := strings.Split(strings.Replace(path, "/", "", 1), "/")
		folder := cache.Dir
		file := ""
		for i, n := range paths {
			if i == (len(paths) - 1) {
				file = n
				os.MkdirAll(folder, 0777)
			} else {
				folder += n + "/"
			}
		}

		filePath := folder + file

		c := cache.Cache{}
		if c.Exists(path) {
			if _, err := os.Stat(filePath); err == nil {
				w.Header().Set("Content-Type", c.ContentType)
				http.ServeFile(w, r, filePath)
				traffic := cache.Traffic{}
				go traffic.Create(strings.Split(r.RemoteAddr, ":")[0], c.ID)
			} else {
				log.Println("Redirecting the user to the ReverseProxy.")
				w.Header().Set("Location", Instance.ProxyServer.URL+path)
				w.WriteHeader(http.StatusFound)
				log.Println("An error occurred where the cache is stored in the database but not in the cache folder.")
				go c.Delete(c.ID)
			}
		} else {
			response, err := http.Get(URL + path)
			Treagler.Do(func() error { return err }, fail, URL+path)

			if response.StatusCode != http.StatusOK {
				defer response.Body.Close()
				log.Println(path, "responded with", response.Status)
				w.WriteHeader(response.StatusCode)
				return
			}

			log.Println("Redirecting the user to the ReverseProxy.")
			w.Header().Set("Location", Instance.ProxyServer.URL+path)
			w.WriteHeader(http.StatusFound)

			go func() {
				if queue.Exists(path) {
					return
				} else {
					queue.Create(path, filePath)
				}

				if _, err := os.Stat(filePath); err == nil {
					err := os.Remove(filePath)
					if err != nil {
						panic(err)
					}
				}

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

				log.Println("Finished downloading:", path, "- Size: ", written/1000000, "MB.")
				queue.Remove(path)
				c.Create(path, filePath, response.Header.Get("Content-Type"))
				traffic := cache.Traffic{}
				traffic.Create(strings.Split(r.RemoteAddr, ":")[0], c.ID)
				if cache.SizeLeft() < int(written) {
					removedCache := cache.SmallestTraffic()
					err := os.Remove(removedCache.File)
					log.Println(
						"Cache",
						removedCache.ID,
						"got removed because of reaching the maximum cache size allowed.",
					)
					removedCache.Delete(removedCache.ID)
					if err != nil {
						panic(err)
					}
				}
			}()
		}
	}
}
