package server_groups

import "github.com/gin-gonic/gin"

type ServerGroupsHandler struct{}

func (*ServerGroupsHandler) CreateServer(cxt *gin.Context) {}

func (*ServerGroupsHandler) UpdateServer(cxt *gin.Context)      {}
func (*ServerGroupsHandler) RemoveServer(cxt *gin.Context)      {}
func (*ServerGroupsHandler) GetServers(cxt *gin.Context)        {}
func (*ServerGroupsHandler) GetServerListPage(cxt *gin.Context) {}
