package auth

import (
	"im-services/internal/api/interfaces"
	"im-services/internal/api/services"
	"im-services/internal/config"
	"im-services/pkg/jwt"
	"im-services/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (*OAuthHandler) GithubOAuth(cxt *gin.Context) {
	var err error
	var code = cxt.Query("code")
	var loginType = cxt.Query("login_type")

	var oauth interfaces.OAuth

	if loginType == "gitee" {
		oauth = new(services.GiteeOAuthService)
	} else {
		oauth = new(services.GithubOAuthService)
	}

	var tokenAuthUrl = oauth.GetTokenAuthUrl(code)
	var token *services.Token

	if token, err = oauth.GetToken(tokenAuthUrl); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(cxt)
		return
	}
	var userInfo map[string]interface{}
	userInfo, err = oauth.GetUserInfo(token.AccessToken)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(cxt)
		return
	}

	err, users, _ := auth.CreateOauthUser(userInfo, loginType)
	ttl := config.Conf.JWT.Ttl
	expireAtTime := time.Now().Unix() + ttl
	tokens := jwt.NewJWT().IssueToken(
		users.ID,
		users.Uid,
		users.Name,
		users.Email,
		expireAtTime,
	)
	//if isNew {
	//	//新注册用户 投递消息
	//	services.InitChatBotMessage(1, users.ID)
	//
	//}

	response.SuccessResponse(&loginResponse{
		ID:         users.ID,
		UID:        users.Uid,
		Name:       users.Name,
		Avatar:     users.Avatar,
		Email:      users.Email,
		ExpireTime: expireAtTime,
		Token:      tokens,
		Ttl:        ttl,
	}).WriteTo(cxt)

}
