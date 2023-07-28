package friend

import (
	"im-services/internal/api/handler"
	"im-services/internal/dao/friend_dao"
	"im-services/internal/enum"
	"im-services/internal/helpers"
	"im-services/internal/service/dispatch"
	"im-services/pkg/response"

	"github.com/gin-gonic/gin"
)

type FriendHandler struct {
}

var (
	friendDao friend_dao.FriendDao
)

func (*FriendHandler) Index(cxt *gin.Context) {
	id := cxt.MustGet("id")

	err, lists := friendDao.GetFriendLists(id)

	if err != nil {
		response.FailResponse(enum.ParamError, "获取用户列表失败").ToJson(cxt)
		return
	}
	response.SuccessResponse(lists).ToJson(cxt)
	return

}

func (*FriendHandler) Show(cxt *gin.Context) {

	err, person := handler.GetPersonId(cxt)
	if err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(cxt)
		return
	}

	var friendDao friend_dao.FriendDao

	err, lists := friendDao.GetFriends(person.ID)

	if err != nil {
		response.SuccessResponse().ToJson(cxt)
		return
	}
	response.SuccessResponse(&lists).ToJson(cxt)
	return
}

func (*FriendHandler) Delete(cxt *gin.Context) {
	err, person := handler.GetPersonId(cxt)
	if err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(cxt)
		return
	}
	var friendDao friend_dao.FriendDao

	errs := friendDao.DelFriends(person.ID, cxt.MustGet("id"))
	if errs != nil {
		response.FailResponse(enum.ParamError, errs.Error()).ToJson(cxt)
		return
	}
	response.SuccessResponse().ToJson(cxt)
	return
}

func (*FriendHandler) GetUserStatus(cxt *gin.Context) {

	err, person := handler.GetPersonId(cxt)
	if err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(cxt)
		return
	}

	var _dispatch dispatch.DispatchService
	ok, _ := _dispatch.IsDispatchNode(person.ID)

	if ok {
		response.SuccessResponse(&UserStatus{
			Status: enum.WsUserOnline,
			Id:     helpers.StringToInt(person.ID),
		}).ToJson(cxt)
		return
	}

	response.SuccessResponse(&UserStatus{
		Status: enum.WsUserOffline,
		Id:     helpers.StringToInt(person.ID),
	}).ToJson(cxt)
	return

}
