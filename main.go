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
	"os"
)

func init() {
	cache.Dir = os.Getenv("ATCACHE_DIR")
	if cache.Dir == "" {
		os.Mkdir("./caches/", 0777)
		cache.Dir = "./caches/"
	} else {
		os.MkdirAll(cache.Dir, 0777)
	}

	if err := cache.LoadDB("root", "", "at_cache"); err != nil {
		panic(err)
	}
}

func main() {
	var port int

	flag.IntVar(&port, "port", 1818, "the port of the server")
	flag.IntVar(&cache.MaxSize, "maxsize", 1000000000 /* 1 GB */, "the max size of the cache")
	flag.Parse()
	server.Instance.Start(fmt.Sprint(port))
	fmt.Println("ATCache is running on port " + fmt.Sprint(port) + ", press ENTER to quit...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	server.Instance.Shutdown()
}
