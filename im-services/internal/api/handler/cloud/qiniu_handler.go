package cloud

import (
	"fmt"
	"im-services/internal/api/services"
	"im-services/internal/config"
	"im-services/internal/enum"
	"im-services/pkg/response"

	"github.com/gin-gonic/gin"
)

type QiNiuHandler struct {
}

var (
	Service services.QiNiuService
)

func (qiniu *QiNiuHandler) UploadFile(cxt *gin.Context) {

	file, err := cxt.FormFile("file")
	if err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(cxt)
		return
	}

	filePath := config.Conf.Server.FilePath + "/" + file.Filename
	err = cxt.SaveUploadedFile(file, filePath)
	if err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(cxt)
		return
	}
	var res Response
	// todo: 图片上传到七牛云
	// fileUrl, _ := Service.UploadFile(filePath, file.Filename)
	// res.FileUrl = fileUrl
	res.FileUrl = filePath
	fmt.Println("res.FileUrl=", res.FileUrl)
	response.SuccessResponse(res).ToJson(cxt)
	return
}

type Response struct {
	FileUrl string `json:"file_url"`
}
