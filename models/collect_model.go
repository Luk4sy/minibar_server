package models

type CollectModel struct {
	Model
	Title        string    `gorm:"size:32"  json:"title"`
	Abstract     string    `gorm:"size:256" json:"abstract"`
	Cover        string    `gorm:"size:256" json:"cover"`
	ArticleCount int       `json:"articleCount"` // 收藏夹中文章数量
	UserID       uint      `json:"userID"`
	IsDefault    bool      `json:"isDefault"` // 是否为默认收藏夹
	UserModel    UserModel `gorm:"foreignKey:UserID" json:"-"`
}
