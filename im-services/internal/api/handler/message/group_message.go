package message

import (
	"im-services/internal/helpers"
	"im-services/internal/models/group_message"
	"im-services/pkg/model"
	"im-services/pkg/response"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type GroupMessageHandler struct {
}

func (*GroupMessageHandler) Index(cxt *gin.Context) {
	page := cxt.Query("page")
	groupId := cxt.Query("to_id")
	pageSize := helpers.StringToInt(cxt.DefaultQuery("pageSize", "50"))

	var list []group_message.ImGroupMessages

	query := model.DB.Model(&group_message.ImGroupMessages{}).Preload("Users").
		Where("group_id=?", groupId).
		Order("send_time desc")

	var total int64
	query.Count(&total)

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

	sortByGroupMessage(list)
	response.SuccessResponse(gin.H{
		"list": list,
		"mate": gin.H{
			"pageSize": pageSize,
			"page":     page,
			"total":    total,
		}}, http.StatusOK).ToJson(cxt)
	return

}

// 对群聊消息进行排序
func sortByGroupMessage(list []group_message.ImGroupMessages) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].SendTime < list[j].SendTime
	})
}
