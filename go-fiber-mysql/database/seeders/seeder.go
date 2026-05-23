package seeders

import (
	"fmt"
	"go-fiber-mysql/internal/app/entities"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	db.Transaction(func(tx *gorm.DB) error {
		userSeeder(tx)

		return nil
	})
}

func userSeeder(db *gorm.DB) {
	datas := []entities.User{
		{
			Name:     "admin",
			Email:    "admin@gmail.com",
			Password: "password",
		},
		{
			Name:     "user",
			Email:    "user@gmail.com",
			Password: "password",
		},
	}

	result := db.Create(&datas)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("%d record of users created\n", result.RowsAffected)
}
