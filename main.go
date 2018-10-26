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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/AnimeTwist/ATCache/cache"
	"github.com/AnimeTwist/ATCache/server"
	"log"
	"os"
)

func init() {
	log.SetPrefix("(ATCache) ")
	cache.Dir = os.Getenv("ATCACHE_DIR")
	if cache.Dir == "" {
		os.Mkdir("./caches/", 0777)
		cache.Dir = "./caches/"
	} else {
		os.MkdirAll(cache.Dir, 0777)
	}
	log.Println("Cache folder path set to ", cache.Dir)

	if err := cache.LoadDB("root", "", "at_cache"); err != nil {
		panic(err)
	}
	log.Println("Loaded the mysql database successfully.")
}

func main() {
	var port int

	flag.IntVar(&port, "port", 1818, "the port of the server")
	flag.StringVar(&server.URL, "url", "http://localhost", "the url of the site to cache")
	flag.IntVar(&cache.MaxSize, "maxsize", 1000000000*1000 /* 1 TB */, "the max size of the cache")
	flag.Parse()

	server.Instance.Start(fmt.Sprint(port))
	log.Println("ATCache is running on port", port, "press ENTER to quit...")

	bufio.NewScanner(os.Stdin).Scan()
	log.Println("Killing processes...")
	server.Instance.Shutdown()
	cache.Instance.Close()
	log.Println("Successfully shutdown everything. Bye, senpai.")
}
