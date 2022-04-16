package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "201911060735",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				gorm.Model
				Name               string `json:"name"`
				FirstName          string `json:"first_name"`
				LastName           string `json:"last_name"`
				Password           string `json:"password" valid:"required"`
				Email              string `json:"email" valid:"email,required"`
				Phone              string `json:"phone"`
				Language           string `json:"language"`
				ActivationKey      string `json:"activation_key"`
				ResetKey           string `json:"reset_key"`
				Active             bool   `json:"active"`
				Admin              bool   `json:"admin"`
				LastModification   int64  `json:"last_access"`
				LastPasswordUpdate int64  `json:"last_password_update"`
			}
			return tx.AutoMigrate(&User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("users").Error
		},
	})
}
