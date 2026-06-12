package seeders

import (
	"fmt"
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/enums"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := permissionSeeder(tx); err != nil {
			return err
		}

		permissions := []entities.Permission{}
		if err := tx.Find(&permissions).Error; err != nil {
			return err
		}

		if err := roleSeeder(tx, permissions); err != nil {
			return err
		}

		if err := userSeeder(tx); err != nil {
			return err
		}

		return nil
	})
}

func permissionSeeder(db *gorm.DB) error {
	var permissions []entities.Permission
	for _, p := range enums.GetAllPermissions() {
		permissions = append(permissions, entities.Permission{
			Name: p.String(),
		})
	}

	for _, p := range permissions {
		var existing entities.Permission
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&p).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	var allPermissions []entities.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	fmt.Printf("%d record of permissions seeded\n", len(allPermissions))
	return nil
}

func roleSeeder(db *gorm.DB, permissions []entities.Permission) error {
	roles := []entities.Role{
		{Name: "User"},
	}

	for _, r := range roles {
		var existing entities.Role
		if err := db.Where("name = ?", r.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&r).Error; err != nil {
					return err
				}
				existing = r
			} else {
				return err
			}
		}

		if existing.Name == "User" {
			// Insert role-permission associations directly
			for _, p := range permissions {
				if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", existing.ID, p.ID).Error; err != nil {
					return err
				}
			}
		}
	}

	fmt.Printf("%d record of roles seeded\n", len(roles))
	return nil
}

func userSeeder(db *gorm.DB) error {
	users := []struct {
		entities.User
		RoleName string
	}{
		{
			User: entities.User{
				Name:     "admin",
				Email:    "admin@gmail.com",
				Password: "password",
				Type:     enums.UserTypeAdmin,
			},
			RoleName: "", // Admin has no roles
		},
		{
			User: entities.User{
				Name:     "user",
				Email:    "user@gmail.com",
				Password: "password",
				Type:     enums.UserTypeUser,
			},
			RoleName: "User",
		},
	}

	for _, u := range users {
		var existing entities.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&u.User).Error; err != nil {
					return err
				}
				existing = u.User
			} else {
				return err
			}
		}

		// Sync Role if RoleName is provided
		if u.RoleName != "" {
			var role entities.Role
			if err := db.Where("name = ?", u.RoleName).First(&role).Error; err != nil {
				return err
			}
			if err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING", existing.ID, role.ID).Error; err != nil {
				return err
			}
		}
	}

	fmt.Printf("%d record of users seeded\n", len(users))
	return nil
}
