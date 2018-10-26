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
	"log"
	"time"
)

var TrafficCache = make(map[uint]map[string]bool)

type Traffic struct {
	Model
	CacheID uint   `gorm:"Type:int(10) unsigned;Column:cache_id;NOT NULL;primary_key" json:"cache_id"`
	Address string `gorm:"Type:varchar(64);Column:address;NOT NULL" json:"name"`
}

func (t *Traffic) Create(addr string, cacheID uint) *Traffic {
	// To prevent from duplicating because of the partial response (206) when streaming the video.
	if _, n := TrafficCache[cacheID]; !n {
		TrafficCache[cacheID] = make(map[string]bool)
	}
	if _, n := TrafficCache[cacheID][addr]; !n {
		t.Address = addr
		t.CacheID = cacheID
		if err := Instance.Create(t).Error; err != nil {
			panic(err)
		}
		TrafficCache[cacheID][addr] = true
		log.Println(addr, "is viewing cache id:", cacheID)
	}

	return t
}

func SmallestTraffic() (c Cache) {
	var caches []Cache

	traffics := make(map[uint]int)
	Instance.Find(&caches)
	for i, cache := range caches {
		if cache.CreatedAt.Unix() < (time.Now().Unix() - int64(time.Hour*24)) {
			for _, traffic := range cache.Traffics {
				if traffic.CreatedAt.Unix() > (time.Now().Unix() - int64(time.Hour*24)) {
					traffics[cache.ID]++
				}
				if i != 0 {
					if traffics[c.ID] > traffics[cache.ID] {
						c = cache
					}
				} else {
					c = cache
				}
			}
		}
	}

	// If all cached files are new then this will prevent from an error to occur
	if len(c.Traffics) == 0 {
		for i, cache := range caches {
			if i != 0 {
				if len(cache.Traffics) < len(c.Traffics) {
					c = cache
				}
			} else {
				c = cache
			}
		}
	}

	return
}
