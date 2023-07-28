package group

import (
	"encoding/json"
	"fmt"
	"im-services/internal/api/handler"
	"im-services/internal/api/requests"
	"im-services/internal/api/services"
	"im-services/internal/dao/group_dao"
	"im-services/internal/enum"
	"im-services/internal/helpers"
	"im-services/internal/models/im_group_users"
	"im-services/internal/models/im_groups"
	"im-services/internal/models/im_messages"
	"im-services/internal/models/im_sessions"
	"im-services/internal/models/user"
	"im-services/pkg/date"
	"im-services/pkg/hash"
	"im-services/pkg/model"
	"im-services/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	groupDao       group_dao.GroupDao
	messageService services.ImMessageService
)

type GroupHandler struct {
}

// 获取群聊列表
func (*GroupHandler) Index(cxt *gin.Context) {
	list := groupDao.GetGroupList(cxt.MustGet("id"))
	response.SuccessResponse(list).ToJson(cxt)
	return
}

func (*GroupHandler) Store(cxt *gin.Context) {
	id := cxt.MustGet("id")
	var selectUser SelectUser

	cxt.ShouldBind(&selectUser)
	selectUser.SelectUser = append(selectUser.SelectUser, helpers.InterfaceToInt64String(id))
	fmt.Println(selectUser.SelectUser)
	params := requests.CreateGroupRequest{
		UserId:     helpers.InterfaceToInt64(id),
		Name:       cxt.PostForm("name"),
		Info:       cxt.PostForm("info"),
		Avatar:     cxt.PostForm("avatar"),
		Password:   cxt.PostForm("password"),
		IsPwd:      helpers.StringToInt(cxt.PostForm("is_pwd")),
		Theme:      cxt.PostForm("theme"),
		SelectUser: selectUser.SelectUser,
	}
	errs := validator.New().Struct(params)
	if errs != nil {
		response.ErrorResponse(enum.ParamError, errs.Error()).ToJson(cxt)
		return
	}
	if params.IsPwd == im_groups.IS_PWD_YES {
		params.Password = hash.BcryptHash(params.Password)
	}
	err, imGroups := groupDao.CreateGroup(params)
	if err != nil {
		response.FailResponse(enum.ApiError, "创建群聊失败！").WriteTo(cxt)
		return
	}

	groupDao.CreateSelectGroupUser(selectUser.SelectUser, int(imGroups.Id), params.Avatar, params.Name)

	// todo 创建成功之后发送创建群聊消息 --

	messageService.SendGroupSessionMessage(selectUser.SelectUser, imGroups.Id)

	response.SuccessResponse(imGroups).WriteTo(cxt)
	return
}

func (*GroupHandler) ApplyJoin(cxt *gin.Context) {
	id := cxt.MustGet("id")
	err, person := handler.GetPersonId(cxt)
	if err != nil {
		response.FailResponse(enum.ParamError, "参数错误！").WriteTo(cxt)
		return
	}
	var group im_groups.ImGroups
	if result := model.DB.Model(&im_groups.ImGroups{}).Where("id=?", person.ID).Find(&group); result.RowsAffected == 0 {
		response.FailResponse(enum.ParamError, "群聊不存在！").WriteTo(cxt)
		return
	}

	if groupDao.IsGroupsUser(id, person.ID) {
		response.FailResponse(enum.ParamError, "已经是群成员了~").WriteTo(cxt)
		return
	}
	if group.IsPwd == int8(im_groups.IS_PWD_YES) {
		if !hash.BcryptCheck(cxt.PostForm("password"), group.Password) {
			response.FailResponse(enum.ParamError, "入群密码错误~,请联系管理员邀请").WriteTo(cxt)
			return
		}
	}

	groupDao.CreateOneGroupUser(group, int(helpers.InterfaceToInt64(id)))

	name := cxt.MustGet("name")

	groupDao.DeleteGroupUser(id, person.ID)

	params := requests.PrivateMessageRequest{
		MsgId:       date.TimeUnixNano(),
		MsgCode:     enum.WsChantMessage,
		MsgClientId: date.TimeUnixNano(),
		FormID:      helpers.InterfaceToInt64(id),
		ToID:        helpers.StringToInt64(person.ID),
		ChannelType: 2,
		MsgType:     im_messages.JOIN_GROUP,
		Message:     fmt.Sprintf("%s 加入群聊", name),
		SendTime:    date.NewDate(),
		Data:        cxt.PostForm("data"),
	}
	// 退群消息推送
	messageService.SendGroupMessage(params)

	response.SuccessResponse().WriteTo(cxt)
	return

	// todo 成功之后 发送入群消息

}

func (*GroupHandler) GetUsers(cxt *gin.Context) {
	err, person := handler.GetPersonId(cxt)
	if err != nil {
		response.FailResponse(enum.ParamError, "参数错误！").WriteTo(cxt)
		return
	}
	var group ImGroups
	if result := model.DB.Model(&im_groups.ImGroups{}).Where("id=?", person.ID).Find(&group); result.RowsAffected == 0 {
		response.FailResponse(enum.ParamError, "群聊不存在！").WriteTo(cxt)
		return
	}
	response.SuccessResponse(&GroupsDate{
		Groups: group,
		Users:  groupDao.GetGroupUsers(person.ID),
	}).WriteTo(cxt)
	return
}

func (*GroupHandler) Logout(cxt *gin.Context) {
	err, person := handler.GetPersonId(cxt)
	if err != nil {
		response.FailResponse(enum.ParamError, "参数错误！").WriteTo(cxt)
		return
	}
	id := cxt.MustGet("id")
	name := cxt.MustGet("name")

	groupDao.DeleteGroupUser(id, person.ID)

	params := requests.PrivateMessageRequest{
		MsgId:       date.TimeUnixNano(),
		MsgCode:     enum.WsChantMessage,
		MsgClientId: date.TimeUnixNano(),
		FormID:      helpers.InterfaceToInt64(id),
		ToID:        helpers.StringToInt64(person.ID),
		ChannelType: im_sessions.GROUP_TYPE,
		MsgType:     im_messages.LOGOUT_GROUP,
		Message:     fmt.Sprintf("%s 退出群聊", name),
		SendTime:    date.NewDate(),
		Data:        cxt.PostForm("data"),
	}
	// 退群消息推送
	messageService.SendGroupMessage(params)

	response.SuccessResponse().WriteTo(cxt)
	return
}

func (*GroupHandler) CreateOrRemoveUser(cxt *gin.Context) {

	var selectUser SelectUser

	cxt.ShouldBind(&selectUser)

	params := requests.CreateUserToGroupRequest{
		GroupId: helpers.StringToInt64(cxt.PostForm("group_id")),
		Type:    helpers.StringToInt(cxt.PostForm("type")),
		UserId:  selectUser.SelectUser,
	}

	userId := cxt.MustGet("id")
	name := cxt.MustGet("name")
	var group ImGroups
	if result := model.DB.Model(&im_groups.ImGroups{}).Where("id=?", params.GroupId).Find(&group); result.RowsAffected == 0 {
		response.FailResponse(enum.ParamError, "群聊不存在！").WriteTo(cxt)
		return
	}
	if group.UserId != userId {
		response.FailResponse(enum.ParamError, "非群主不可以邀请人入群！").WriteTo(cxt)
		return

	}

	if params.Type == 1 {
		groupDao.CreateSelectGroupUser(selectUser.SelectUser, int(params.GroupId), group.Avatar, group.Name)
		// 发送群聊会话消息
		messageService.SendGroupSessionMessage(selectUser.SelectUser, params.GroupId)
	} else {
		groupDao.DelSelectGroupUser(selectUser.SelectUser, int(params.GroupId), group.Avatar, group.Name)
	}
	var users []user.ImUsers

	model.DB.Model(&user.ImUsers{}).
		Where("id in(?)", model.DB.Model(&im_group_users.ImGroupUsers{}).
			Where("group_id=?", params.GroupId).Select("user_id")).
		Find(&users)

	groupStr, _ := json.Marshal(group)
	message := requests.PrivateMessageRequest{
		MsgId:       date.TimeUnixNano(),
		MsgCode:     enum.WsChantMessage,
		MsgClientId: date.TimeUnixNano(),
		FormID:      group.Id,
		ChannelType: enum.GroupMessage,
		MsgType:     enum.JOIN_GROUP,
		Message:     "",
		SendTime:    date.NewDate(),
		Data:        string(groupStr),
	}

	messageService.SendCreateUserGroupMessage(users, message, name, params.Type, selectUser.SelectUser)

	response.SuccessResponse().WriteTo(cxt)
	return

}
