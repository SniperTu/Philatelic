package interfaces

type OAuth interface {
	GetTokenAuthUrl(code string) string
	GetUserInfo(token string) (map[string]interface{}, error)
}
