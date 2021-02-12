package config

type DBConf struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     uint32
}

func DBConfig() DBConf {
	return DBConf{
		Host:     "localhost",
		Port:     5432,
		User:     "wallester",
		Password: "12345",
		DBName:   "wallester",
	}
}
