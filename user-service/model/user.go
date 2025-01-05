package model

type User struct {
	Id        int64  `gorm:"column:id;type:int;primaryKey;autoIncrement;not null;comment:id" json:"id"`
	Name      string `gorm:"column:name;type:varchar(50);not null;comment:用户名" json:"name"`
	Email     string `gorm:"column:email;type:varchar(50);not null;comment:邮箱" json:"email"`
	Password  string `gorm:"column:password;type:varchar(50);not null;comment:密码" json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(20);not null;comment:手机号" json:"phone"`
	Address   string `gorm:"column:address;type:varchar(255);not null;comment:地址" json:"address"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;comment:创建时间" json:"createdAt"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;comment:修改时间" json:"updatedAt"`
}

func (User) TableName() string {
	return "t_user"
}
