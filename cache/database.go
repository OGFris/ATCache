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
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Instance *gorm.DB

func LoadDB(user, password, name string) (err error) {
	Instance, err = gorm.Open("mysql", user+":"+password+"@/"+name+"?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		Instance.AutoMigrate(
			&Cache{},
			&Traffic{},
			&DownloadQueue{},
		)
		Instance.Model(&Traffic{}).AddForeignKey("cache_id", "caches(id)", "RESTRICT", "RESTRICT")
		return nil
	}
	return
}
