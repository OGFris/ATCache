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
	"os"
	"time"
)

var (
	Dir          string
	MaxSize      int
	TrafficCache = make(map[uint]map[string]bool)
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"Default:null" sql:"index" json:"deleted_at"`
}

type Cache struct {
	Model
	Path        string    `gorm:"Type:varchar(64);Column:path;NOT NULL;primary_key;unique" json:"path"`
	File        string    `gorm:"Type:varchar(64);Column:file;NOT NULL;primary_key;unique" json:"file"`
	ContentType string    `gorm:"Type:varchar(64);Column:content_type;NOT NULL" json:"content_type"`
	Traffics    []Traffic `json:"-"`
}

type Traffic struct {
	Model
	CacheID uint   `gorm:"Type:int(10) unsigned;Column:cache_id;NOT NULL;primary_key" json:"cache_id"`
	Address string `gorm:"Type:varchar(64);Column:address;NOT NULL" json:"name"`
}

func (c *Cache) Create(path, file, contentType string) *Cache {
	c.Path = path
	c.File = file
	c.ContentType = contentType
	Instance.Create(c)

	return c
}

func (c *Cache) Delete(id uint) *Cache {
	Instance.Delete(c, c)

	return c
}

func (c *Cache) Exists(path string) bool {
	return Instance.Find(c, &Cache{Path: path}).Error == nil
}

func (t *Traffic) Create(addr string, cacheID uint) *Traffic {
	// To prevent from duplicating because of the partial response (206) when streaming the video.
	if _, n := TrafficCache[cacheID]; !n {
		TrafficCache[cacheID] = make(map[string]bool)
	}
	if _, n := TrafficCache[cacheID][addr]; !n {
		t.Address = addr
		t.CacheID = cacheID
		Instance.Create(t)
		TrafficCache[cacheID][addr] = true
		log.Println(addr, " is viewing cache id: ", cacheID)
	}

	return t
}

func FolderSize() int {
	return folderSize(Dir)
}

func folderSize(path string) (size int) {
	folder, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	files, err := folder.Readdir(0)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			size += folderSize(path + f.Name())
		} else {
			size += int(f.Size())
		}
	}

	return
}

func SizeLeft() int {
	return MaxSize - FolderSize()
}

func SmallestTraffic() (c Cache) {
	var (
		caches []Cache
		last   uint
	)

	traffics := make(map[uint]int)
	Instance.Find(&caches)
	for i, cache := range caches {
		for _, traffic := range cache.Traffics {
			if traffic.CreatedAt.Unix() > (time.Now().Unix() - int64(time.Hour*24*7)) {
				traffics[cache.ID]++
			}
			if i != 0 {
				if traffics[last] > traffics[cache.ID] {
					c = cache
				}
			} else {
				last = cache.ID
			}
		}
	}

	return
}
