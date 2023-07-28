package sever_group

const (
	INVITE_NEED_AGREE     = 0 // 需要被邀请人的同意(默认)
	INVITE_NEED_NOT_AGREE = 1 // 邀请不需要同意
)
const (
	APPLY_NEED_NOT_AGREE = 0 // 申请加入不需要创建者或管理员的同意(默认)
	APPLY_NEED_AGREE     = 1 // 申请需要同意
)

const (
	CAN_SEARCH     = 0 // 可被搜索到（默认）
	CAN_NOT_SEARCH = 1 // 不可被搜索到
)

type ServerGroup struct {
	Id           string `gorm:"column:id" json:"id"`                     // 圈组id
	OwnerId      int64  `json:"owner_id"`                                // 圈组操纵者id
	Name         string `gorm:"column:name" json:"name"`                 // 圈组名称
	Icon         string `gorm:"column:icon" json:"icon"`                 // 圈组图标
	Custom       string `gorm:"column:custom" json:"custom"`             // 自定义扩展
	InviteMode   int8   `gorm:"column:inviteMode" json:"inviteMode"`     // 圈组邀请模式 0，邀请需要被邀请人的同意(默认)；1，邀请不需要同意
	ApplyMode    int8   `gorm:"column:applyMode" json:"applyMode"`       // 圈组申请模式 0，申请加入不需要创建者或管理员的同意(默认)；1，申请需要同意
	SearchEnable int8   `gorm:"column:searchEnable" json:"searchEnable"` //是否可被搜索到 0，否；1，是（默认）
	CreatedAt    string `gorm:"column:created_at" json:"created_at"`     //创建时间
	UpdateAt     string `gorm:"column:update_at" json:"update_at"`       // 更新时间
	DeleteAt     string `gorm:"column:delete_at" json:"delete_at"`       // 删除时间
}
