package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

var fileTemplate = `package migrations
import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)
func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "%s",
		Migrate: func(tx *gorm.DB) error {
			type CHANGEME struct {
				gorm.Model
			}
			return tx.AutoMigrate(&CHANGEME{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("changeme").Error
		},
	})
}
`

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("must exactly one arg (hint_name of migration file) defined, args are now %v.", args)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	hintName := args[1]                                // entity to migrate
	nowTimeString := time.Now().Format("200601021504") // date
	fileName := fmt.Sprintf("migrations/%s_%s.go", nowTimeString, hintName)
	longFileName := path.Join(workingDirectory, fileName)

	fd, err := os.OpenFile(longFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Open file %s failed: %v", fileName, err)
	}

	fileContent := fmt.Sprintf(fileTemplate, nowTimeString)
	_, err = fd.WriteString(fileContent)
	if err != nil {
		log.Fatalf("write file %s failed: %v", fileName, err)
	}

	fmt.Printf("Successfully created migration %s\n", fileName)
}