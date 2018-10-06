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

// TODO:
// [+] Http server to handle the request
// [+] Downloading the videos then storing them
// [ ] Caching the videos with a timer


func init() {
	cache.CacheDir = os.Getenv("AT_CACHE_DIR")

	if cache.CacheDir == "" {
		os.Mkdir("./caches/", 0777)
		cache.CacheDir = "./caches/"
	}


}

func main() {
	var port int

	flag.IntVar(&port, "port", 1818, "set the port of the server")

	server.Instance.Start(fmt.Sprint(port))

	fmt.Println("ATCache is running on port " + fmt.Sprint(port) + ", press ENTER to quit...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	server.Instance.Shutdown()
}
