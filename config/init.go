package config

import (
	"github.com/Godzizizilla/Management-System/cache"
	"github.com/Godzizizilla/Management-System/db"
)

func Init() {
	db.SetupDB()
	cache.SetupRedis()
	InitAdmin()
}
