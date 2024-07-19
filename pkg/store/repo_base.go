package repo

import (
	"github.com/jinzhu/gorm"
	"sync"
)

var _dbIns *BaseDB
var _dbInsOnce sync.Once

type BaseDB struct {
	Engine *gorm.DB
}

func GetBaseDB() *BaseDB {
	_dbInsOnce.Do(func() {
		_dbIns = &BaseDB{
			Engine: InitDB(),
		}
	})
	return _dbIns
}
