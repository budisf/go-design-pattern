package driver

import (
	config "ethical-be/app/config"
	database "ethical-be/app/databases"
)

var (
	/*
	 Define .env configuratian
	*/
	Conf, ErrConf = config.Init()

	/*
	 Database connection
	*/
	DB, _ = database.Connect(&Conf)
)
