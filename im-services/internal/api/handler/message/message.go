package message

import (
	"im-services/internal/api/requests"
	"im-services/internal/api/services"
	"im-services/internal/dao/friend_dao"
	"im-services/internal/dao/messsage_dao"
	"im-services/internal/enum"
	"im-services/internal/helpers"
	"im-services/internal/models/im_messages"
	"im-services/internal/models/user"
	"im-services/pkg/date"
	"im-services/pkg/model"
	"im-services/pkg/response"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MessageHandler struct {
}

var (
	messagesServices services.ImMessageService
	friend           friend_dao.FriendDao
	messageDao       messsage_dao.MessageDao
)

func (m *MessageHandler) Index(cxt *gin.Context) {

	id := cxt.MustGet("id")
	page := cxt.Query("page")
	toId := cxt.Query("to_id")
	pageSize := helpers.StringToInt(cxt.DefaultQuery("pageSize", "50"))

	var list []im_messages.ImMessages

	query := model.DB.Table("im_messages").
		Where("(form_id=? and to_id=?) or (form_id=? and to_id=?)", id, toId, toId, id).
		Order("created_at desc")

	var total int64
	query.Count(&total)

	var users user.ImUsers

	model.DB.Table("im_users").Where("id=?", toId).First(&users)

	if len(page) > 0 {
		query = query.Where("id<?", page)
	}

	if result := query.Limit(pageSize).Find(&list); result.RowsAffected == 0 {
		response.SuccessResponse(gin.H{
			"list": struct {
			}{},
			"mate": gin.H{
				"pageSize": pageSize,
				"page":     page,
				"total":    0,
			}}, http.StatusOK).ToJson(cxt)
		return
	}

	sortByMessage(list, users)
	response.SuccessResponse(gin.H{
		"list": list,
		"mate": gin.H{
			"pageSize": pageSize,
			"page":     page,
			"total":    total,
		}}, http.StatusOK).ToJson(cxt)
	return

}
func sortByMessage(list []im_messages.ImMessages, users user.ImUsers) {
	sort.Slice(list, func(i, j int) bool {
		list[i].Users.ID = users.ID
		list[i].Users.Name = users.Name
		list[i].Users.Email = users.Email
		list[i].Users.Avatar = users.Avatar
		return list[i].Id < list[j].Id
	})
}

func (m *MessageHandler) RecallMessage(cxt *gin.Context) {
	response.SuccessResponse().ToJson(cxt)
	return
}

func (m *MessageHandler) SendVideoMessage(cxt *gin.Context) {

	id := cxt.MustGet("id")
	toId := cxt.PostForm("to_id")

	if !friend.IsFriends(id, toId) {
		response.FailResponse(enum.WsNotFriend, "非好友关系,不能聊天...").ToJson(cxt)
		return
	}
	var users user.ImUsers
	model.DB.Table("im_users").Where("id=?", id).First(&users)

	params := requests.VideoMessageRequest{
		MsgCode:  enum.VideoChantMessage,
		FormID:   helpers.InterfaceToInt64(id),
		ToID:     helpers.StringToInt64(toId),
		Message:  "视频请求...",
		SendTime: date.NewDate(),
		Users: requests.Users{
			Email:  users.Email,
			Name:   users.Name,
			Avatar: users.Avatar,
		},
	}
	ok := messagesServices.SendVideoMessage(params)
	if !ok {
		response.FailResponse(http.StatusInternalServerError, "用户不在线").ToJson(cxt)
		return
	}
	response.SuccessResponse(params).ToJson(cxt)
	return

}

func (m *MessageHandler) SendMessage(cxt *gin.Context) {

	id := cxt.MustGet("id")
	params := requests.PrivateMessageRequest{
		MsgId:       date.TimeUnixNano(),
		MsgCode:     enum.WsChantMessage,
		MsgClientId: helpers.StringToInt64(cxt.PostForm("msg_client_id")),
		FormID:      helpers.InterfaceToInt64(id),
		ToID:        helpers.StringToInt64(cxt.PostForm("to_id")),
		ChannelType: helpers.StringToInt(cxt.DefaultPostForm("channel_type", "1")),
		MsgType:     helpers.StringToInt(cxt.PostForm("msg_type")),
		Message:     cxt.PostForm("message"),
		SendTime:    date.NewDate(),
		Data:        cxt.PostForm("data"),
	}

	errs := validator.New().Struct(params)

	if errs != nil {
		response.FailResponse(enum.ParamError, errs.Error()).ToJson(cxt)
		return
	}

	switch params.ChannelType {
	case 1:
		messageDao.CreateMessage(params)
		// todo 暂时先写死 --
		var users user.ImUsers
		model.DB.Model(&user.ImUsers{}).Where("id =?", params.ToID).Find(&users)
		if users.UserType == user.BOT_TYPE {
			// todo 消息投递 机器人不需要好友关系 在线不发消息 离线发送消息
			if messagesServices.IsOline(helpers.Int64ToString(users.ID)) {
				messagesServices.SendPrivateMessage(params)
			} else {
				messagesServices.SendChatMessage(params)
			}
			response.SuccessResponse(params).ToJson(cxt)
			return
		} else {
			if !friend.IsFriends(id, params.ToID) {
				response.FailResponse(enum.WsNotFriend, "非好友关系,不能聊天...").ToJson(cxt)
				return
			}
			// todo 此处有点逻辑bug
			// 消息投递
			ok, msg := messagesServices.SendPrivateMessage(params)
			if !ok {
				response.FailResponse(http.StatusInternalServerError, msg).ToJson(cxt)
				return
			}
		}
	case 2:
		if !groupDao.IsGroupsUser(id, params.ToID) {
			response.FailResponse(enum.WsNotFriend, "你不是此群成员了...").ToJson(cxt)
			return
		}

		// 消息投递
		ok := messagesServices.SendGroupMessage(params)
		if !ok {
			response.FailResponse(http.StatusInternalServerError, "群聊消息投递异常").ToJson(cxt)
			return
		}

	}

	response.SuccessResponse(params).ToJson(cxt)
	return
}
