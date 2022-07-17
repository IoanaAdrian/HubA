package services

import (
	"github.com/IoanaAdrian/HubA/src/api/services/database"
)

func Init(){
	EnvironmentInit()
	database.DbInit()
}