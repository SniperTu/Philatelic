# Philatelic 后端设计文档

### HTTP路由列表

| 请求方式 |             路由路径             |        模块功能        |
| :------: | :------------------------------: | :--------------------: |
|   POST   |              /login              |         *登录*         |
|   POST   |           /registered            |         *注册*         |
|   POST   |          /sendEmailCode          |     *发送注册邮件*     |
|   GET    |            /user/:id             |     *获取用户信息*     |
|   Any    |          /address/list           |      *通讯录列表*      |
|   GET    |            /sessions             |     *获取会话列表*     |
|   POST   |            /sessions             |       *添加会话*       |
|   PUT    |          /sessions/:id           |       *更新会话*       |
|  DELETE  |          /sessions/:id           |       *移除会话*       |
|   Any    |             /friends             |     *获取好友列表*     |
|   GET    |           /friends/:id           |   *获取好友详情信息*   |
|  DELETE  |           /friends/:id           |       *删除好友*       |
|   GET    |       /friends/status/:id        |     *获取用户状态*     |
|   POST   |         /friends/record          |     *发送好友请求*     |
|   GET    |         /friends/record          | *获取好友申请记录列表* |
|   PUT    |         /friends/record          |     *同意好友请求*     |
|   GET    |        /friends/userQuery        |    *非好友用户查询*    |
|   GET    |            /messages             |   *获取私聊消息列表*   |
|   GET    |         /messages/groups         |   *获取群聊消息列表*   |
|   POST   |        /messages/private         |     *发送私聊消息*     |
|   POST   |         /messages/group          |     *发送群聊消息*     |
|   POST   |         /messages/video          |     *发送视频请求*     |
|   POST   |         /messages/recall         |       *消息撤回*       |
|   POST   |          /groups/store           |       *创建群组*       |
|   POST   |      /groups/applyJoin/:id       |       *加入群组*       |
|   POST   |    /groups/createOrRemoveUser    |     *添加移除用户*     |
|   GET    |           /groups/list           |     *获取群组列表*     |
|   GET    |        /groups/users/:id         |    *获取群成员信息*    |
|  DELETE  |           /groups/:id            |       *退出群聊*       |
|   POST   |           /invite/:id            |  *创建分享群聊token*   |
|   POST   |           /upload/file           |       *上传文件*       |
|   POST   |   /server_groups/createServer    |    *创建圈组服务器*    |
|   POST   |   /server_groups/updateServer    |  *修改圈组服务器信息*  |
|   POST   |   /server_groups/removeServer    |    *删除圈组服务器*    |
|   POST   |    /server_groups/getServers     |  *批量查询服务器信息*  |
|   POST   | /server_groups/getServerListPage |  *分页查询服务器列表*  |
|          |                                  |                        |



### 数据表设计

### 



