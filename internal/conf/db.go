package conf

import "fmt"

func MySQLDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Database.User,
		Database.Password,
		Database.Host,
		Database.Port,
		Database.Name,
	)
}

func PostgresDsn() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		Database.Host,
		Database.Port,
		Database.User,
		Database.Password,
		Database.Name,
	)
	if Database.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", Database.Schema)
	}
	return dsn
}
