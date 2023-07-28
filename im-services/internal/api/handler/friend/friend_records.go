package friend

import (
	"im-services/internal/api/requests"
	"im-services/internal/api/services"
	"im-services/internal/dao/friend_dao"
	"im-services/internal/dao/session_dao"
	"im-services/internal/enum"
	"im-services/internal/helpers"
	"im-services/internal/models/im_friend_records"
	"im-services/internal/models/im_friends"
	"im-services/internal/models/user"
	"im-services/internal/service/client"
	"im-services/pkg/date"
	"im-services/pkg/model"
	"im-services/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FriendRecordHandler struct {
}

func (friend *FriendRecordHandler) Index(cxt *gin.Context) {
	var list []im_friend_records.ImFriendRecords
	id := cxt.MustGet("id")
	if result := model.DB.Model(&im_friend_records.ImFriendRecords{}).Preload("Users").
		Where("to_id=? or form_id=?", id, id).
		Order("created_at desc").Find(&list); result.RowsAffected == 0 {
		response.SuccessResponse().ToJson(cxt)
		return
	}

	response.SuccessResponse(list).ToJson(cxt)
	return

}

func (friend *FriendRecordHandler) Store(cxt *gin.Context) {
	id := cxt.MustGet("id")

	params := requests.CreateFriendRequest{
		ToId:        cxt.PostForm("to_id"),
		Information: cxt.PostForm("information"),
	}

	errs := validator.New().Struct(params)

	if errs != nil {
		response.ErrorResponse(enum.ParamError, errs.Error()).ToJson(cxt)
		return
	}
	var users user.ImUsers

	if result := model.DB.Table("im_users").Where("id=?", params.ToId).First(&users); result.RowsAffected == 0 {
		response.ErrorResponse(enum.ParamError, "用户不存在").ToJson(cxt)
		return
	}

	var record im_friend_records.ImFriendRecords
	if result := model.DB.Table("im_friend_records").Where("form_id=? and to_id=? and status=0", id, params.ToId).First(&record); result.RowsAffected > 0 {
		response.ErrorResponse(enum.ParamError, "请勿重复添加...").ToJson(cxt)
		return
	}

	var friends im_friends.ImFriends

	if result := model.DB.Table("im_friends").Where("form_id=? and to_id=?", id, params.ToId).First(&friends); result.RowsAffected > 0 {
		response.ErrorResponse(enum.ParamError, "用户已经是好友关系了...").ToJson(cxt)
		return
	}

	records := im_friend_records.ImFriendRecords{
		FormId:      helpers.InterfaceToInt64(id),
		ToId:        helpers.StringToInt64(params.ToId),
		Status:      im_friend_records.WaitingStatus,
		CreatedAt:   date.NewDate(),
		Information: params.Information,
	}

	model.DB.Save(&records)

	var messageService services.ImMessageService

	var msg client.CreateFriendMessage

	msg.MsgCode = enum.WsCreate
	msg.ID = records.Id
	msg.ToID = records.ToId
	msg.FormId = records.FormId
	msg.Information = records.Information
	msg.CreatedAt = records.CreatedAt
	msg.Status = records.Status
	msg.Users.ID = users.ID
	msg.Users.Avatar = users.Avatar
	msg.Users.Name = users.Name

	messageService.SendFriendActionMessage(msg)

	records.Users.Name = users.Name
	records.Users.Id = users.ID
	records.Users.Avatar = users.Avatar
	response.SuccessResponse(records).ToJson(cxt)
	return
}

func (friend *FriendRecordHandler) Update(cxt *gin.Context) {
	id := cxt.MustGet("id")
	params := requests.UpdateFriendRequest{
		Status: helpers.StringToInt(cxt.PostForm("status")),
		ID:     cxt.PostForm("id"),
	}
	errs := validator.New().Struct(params)
	if errs != nil {
		response.ErrorResponse(enum.ParamError, errs.Error()).ToJson(cxt)
		return
	}
	var records im_friend_records.ImFriendRecords

	if result := model.DB.Table("im_friend_records").
		Where("id=? and status=0", params.ID).First(&records); result.RowsAffected == 0 {
		response.ErrorResponse(http.StatusInternalServerError, "数据不存在").ToJson(cxt)
		return
	}

	var friends im_friends.ImFriends

	if result := model.DB.Table("im_friends").Where("form_id=? and to_id=?", records.ToId, id).First(&friends); result.RowsAffected > 0 {
		response.ErrorResponse(enum.ParamError, "用户已经是好友关系了...").ToJson(cxt)
		return
	}

	var users user.ImUsers

	model.DB.Table("im_users").Where("id=?", id).Find(&users)

	records.Status = params.Status

	model.DB.Updates(&records)

	var messageService services.ImMessageService

	var msg client.CreateFriendMessage
	var msgCode int

	if params.Status == 1 {
		msgCode = enum.WsFriendOk
		var friendDao friend_dao.FriendDao
		//添加好友关系
		friendDao.AgreeFriendRequest(records.FormId, records.ToId)
		friendDao.AgreeFriendRequest(records.ToId, records.FormId)

		// 添加会话关系
		var sessionDao session_dao.SessionDao
		sessionDao.CreateSession(records.FormId, records.ToId, 1)
		sessionDao.CreateSession(records.ToId, records.FormId, 1)

	} else {
		msgCode = enum.WsFriendError
	}

	msg.MsgCode = msgCode
	msg.ID = records.Id
	msg.ToID = records.FormId
	msg.FormId = records.ToId
	msg.Information = records.Information
	msg.CreatedAt = records.CreatedAt
	msg.Status = records.Status
	msg.Users.ID = users.ID
	msg.Users.Avatar = users.Avatar
	msg.Users.Name = users.Name

	messageService.SendFriendActionMessage(msg)
	friends.Status = params.Status
	response.SuccessResponse(friends).WriteTo(cxt)
	return

}

// QueryUser 查询非好友用户
func (friend *FriendRecordHandler) UserQuery(cxt *gin.Context) {
	id := cxt.MustGet("id")
	params := requests.QueryUserRequest{
		Email: cxt.Query("email"),
	}
	errs := validator.New().Struct(params)
	if errs != nil {
		response.ErrorResponse(enum.ParamError, errs.Error()).ToJson(cxt)
		return
	}
	var friendDao friend_dao.FriendDao
	users := friendDao.GetNotFriendList(id, params.Email)
	response.SuccessResponse(users).ToJson(cxt)
	return
}
