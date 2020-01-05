package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202001051228",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				gorm.Model
				Name               string `json:"name" gorm:"type:text"`
				FirstName          string `json:"first_name" gorm:"type:text"`
				LastName           string `json:"last_name" gorm:"type:text"`
				Password           string `json:"password" gorm:"type:text"`
				Email              string `json:"email" gorm:"type:varchar(100);unique"`
				Role               string `json:"role" gorm:"type:text"`
				Address            string `json:"address" gorm:"type:text"`
				Status             string `json:"status" gorm:"type:text"`
				Phone              string `json:"phone" gorm:"type:text"`
				Language           string `json:"language" gorm:"type:text"`
				ActivationKey      string `json:"activation_key" gorm:"type:text"`
				ResetKey           string `json:"reset_key" gorm:"type:text"`
				LastModification   int64  `json:"last_modification" `
				LastPasswordUpdate int64  `json:"last_password_update" `
				GroupId            string `json:"group_id" gorm:"type:text"`
			}
			return tx.AutoMigrate(&User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("users").Error
		},
	})
}
