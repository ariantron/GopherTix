package models

type Role struct {
	ID   uint8  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50)" json:"name"`
}
