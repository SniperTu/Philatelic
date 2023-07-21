package model

import (
	"fmt"
	"im-services/internal/config"
	"im-services/internal/models/group_invites"
	"im-services/internal/models/group_message"
	"im-services/internal/models/im_friend_records"
	"im-services/internal/models/im_friends"
	"im-services/internal/models/im_group_users"
	"im-services/internal/models/im_groups"
	"im-services/internal/models/im_messages"
	"im-services/internal/models/im_sessions"
	"im-services/internal/models/offline_message"
	"im-services/internal/models/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type BaseModel struct {
	ID int64
}

func InitDb() *gorm.DB {
	var (
		host     = config.Conf.Mysql.Host
		port     = config.Conf.Mysql.Port
		database = config.Conf.Mysql.Database
		username = config.Conf.Mysql.Username
		password = config.Conf.Mysql.Password
		charset  = config.Conf.Mysql.Charset

		err error
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		username, password, host, port, database, charset)

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("Mysql 连接异常: ")
		panic(err)
	}
	fmt.Println("Mysql 连接成功！！！！！")
	err = DB.AutoMigrate(
		&group_invites.ImGroupInvites{},
		&group_message.ImGroupMessages{},
		&im_friend_records.ImFriendRecords{},
		&im_friends.ImFriends{},
		&im_group_users.ImGroupUsers{},
		&im_groups.ImGroups{},
		&im_messages.ImMessages{},
		&im_sessions.ImSessions{},
		&offline_message.ImOfflineMessages{},
		&offline_message.ImGroupOfflineMessages{},
		&user.ImUsers{},
	)
	if err != nil {
		panic("Failed to create tables")
	}
	return DB
}
