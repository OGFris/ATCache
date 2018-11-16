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

package queue

import (
	"github.com/AnimeTwist/ATCache/redis"
	"time"
)

func Exists(path string) bool {
	if redis.Client.Get(path).Err() != nil {
		return false
	}
	return true
}

func Create(path, file string) {
	go func() {
		if err := redis.Client.Set(path, file, time.Second*5).Err(); err != nil {
			panic(err)
		}
	}()
}

func Remove(path string) {
	go func() {
		if err := redis.Client.Del(path).Err(); err != nil {
			panic(err)
		}
	}()
}
