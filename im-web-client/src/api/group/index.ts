import request from '../request'
import type {
  createGroupDataType,
  groupUserType,
  groupInfoType,
  removeType,
} from './type'
// 获取群聊用户信息
export function getGroupUserInfo(params: { id: number }) {
  return request.get<groupInfoType<groupUserType>>(`/groups/users/${params.id}`)
}
//创建群聊
export function createGroup(data: createGroupDataType) {
  return request.post('/groups/store', data)
}

//申请加入群聊
export function applyAddGroup(data: { id: number }) {
  return request.post(`/groups/applyJoin/${data.id}`)
}
// 退出群聊
export function deleteGroup(params: { id: number }) {
  return request.delete(`/groups/${params.id}`)
}
// 邀请或者移除用户入群
export function createOrRemoveUser(data: removeType) {
  return request.post('/groups/createOrRemoveUser', data)
}

// todo: 获取群组列表
export function groupsList() {
  return request.post('/groups/list')
}

// todo: 获取群成员信息
export function getGroupUsers(data: { id: number }) {
  return request.post(`/groups/users/${data.id}`)
}

// todo: 创建分享群聊token
export function invite(data: { id: number }) {
  return request.post(`/invite/${data.id}`)
}