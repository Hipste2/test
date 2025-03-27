package global

import (
	"gorm.io/gorm"
	"mxshop/user_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
)

//func init() {
//	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
//	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//
//	// 设置全局logger，可以在每次执行sql语句的时候会打印每一行sql
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//		logger.Config{
//			SlowThreshold:             time.Second, // Slow SQL threshold
//			LogLevel:                  logger.Info, // Log level
//			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
//			ParameterizedQueries:      true,        // Don't include params in the SQL log
//			Colorful:                  true,        // Disable color
//		},
//	)
//
//	//全局模式
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic(err)
//	}
//}
