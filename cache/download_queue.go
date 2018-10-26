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

type DownloadQueue struct {
	Model
	Path string `gorm:"Type:varchar(64);Column:path;NOT NULL;primary_key;unique" json:"path"`
	File string `gorm:"Type:varchar(64);Column:file;NOT NULL;primary_key;unique" json:"file"`
}

func (d *DownloadQueue) Exists(path string) bool {
	return Instance.Find(d, &DownloadQueue{Path: path}).Error == nil
}

func (d *DownloadQueue) Create(path, file string) *DownloadQueue {
	d.File = file
	d.Path = path
	if err := Instance.Create(d).Error; err != nil {
		panic(err)
	}

	return d
}

func (d *DownloadQueue) Delete(path string) *DownloadQueue {
	if err := Instance.Delete(d, &DownloadQueue{Path: path}).Error; err != nil {
		panic(err)
	}

	return d
}
