package dao

import (
	"im-services/internal/models/offline_message"
	"im-services/pkg/date"
	"im-services/pkg/model"

	"github.com/golang-module/carbon"
)

// PullPrivateOfflineMessage 拉取离线私聊消息
func (offline *OfflineMessageDao) PullPrivateOfflineMessage(id string) []offline_message.ImOfflineMessages {

	var list []offline_message.ImOfflineMessages

	// 拉去最近半个月内的消息记录
	timeStamp := carbon.Parse(date.NewDate()).SubDays(15).Timestamp()

	model.DB.Table("im_offline_messages").
		Where("status=0 and receive_id=? and send_time>?", id, timeStamp).
		Find(&list)

	return list
}

func (offline *OfflineMessageDao) PullPrivateGroupOfflineMessage(id string) []offline_message.ImGroupOfflineMessages {

	var list []offline_message.ImGroupOfflineMessages

	// 拉去最近半个月内的消息记录
	timeStamp := carbon.Parse(date.NewDate()).SubDays(15).Timestamp()

	model.DB.Model(&offline_message.ImGroupOfflineMessages{}).
		Where("status=0 and receive_id=? and send_time>?", id, timeStamp).
		Find(&list)

	return list
}

// UpdatePrivateOfflineMessageStatus 更新消息状态
func (offline *OfflineMessageDao) UpdatePrivateOfflineMessageStatus(id string, ChannelType int) {
	if ChannelType == 1 {
		model.DB.Model(&offline_message.ImOfflineMessages{}).
			Where("status=0 and receive_id=?", id).
			Updates(map[string]interface{}{"status": 1})
	} else {
		model.DB.Model(&offline_message.ImGroupOfflineMessages{}).
			Where("status=0 and receive_id=?", id).
			Updates(map[string]interface{}{"status": 1})
	}
}
