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

package cache

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"time"
)

type Cache struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Expire  time.Time `json:"expire"`
	Traffic uint      `json:"traffic"`
}

type Caches struct {
	file   string
	caches map[string]Cache
}

func (c *Caches) UpdateCache(cache Cache) {
	c.caches[cache.Name] = cache
	c.Save()
}

func (c *Caches) AddCache(cache Cache) {
	c.UpdateCache(cache)
}

func (c *Caches) DeleteCache(cache Cache) {
	delete(c.caches, cache.Name)
	c.Save()
}

func (c *Caches) Save() {
	b, err := jsoniter.Marshal(c.caches)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(c.file, b, 0644)
	if err != nil {
		panic(err)
	}
}
