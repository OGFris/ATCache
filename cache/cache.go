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
	"os"
	"time"
)

var (
	Dir     string
	MaxSize int
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

func (c *Cache) Create(path, file, contentType string) *Cache {
	c.Path = path
	c.File = file
	c.ContentType = contentType
	if err := Instance.Create(c).Error; err != nil {
		panic(err)
	}

	return c
}

func (c *Cache) Delete(id uint) *Cache {
	if err := Instance.Delete(c, c).Error; err != nil {
		panic(err)
	}

	return c
}

func (c *Cache) Exists(path string) bool {
	return Instance.Find(c, &Cache{Path: path}).Error == nil
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
