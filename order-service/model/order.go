package model

type Order struct {
	Id          int64  `gorm:"column:id;type:int;primaryKey;autoIncrement;not null;comment:id" json:"id"`
	UserId      int64  `gorm:"column:user_id;type:bigint;not null;comment:用户ID" json:"userId"`
	Name        string `gorm:"column:name;type:varchar(255);not null;comment:商品名称" json:"name"`
	Price       int64  `gorm:"column:price;type:bigint;not null;comment:订单金额" json:"price"`
	Description string `gorm:"column:description;type:varchar(2550);comment:订单描述" json:"description"`
	CreatedAt   int64  `gorm:"column:created_at;type:bigint;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt   int64  `gorm:"column:updated_at;type:bigint;comment:修改时间" json:"updatedAt"`
}

func (Order) TableName() string {
	return "t_order"
}
