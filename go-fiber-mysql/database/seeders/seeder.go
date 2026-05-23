package seeders

import (
	"fmt"
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/enums"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	db.Transaction(func(tx *gorm.DB) error {
		permissions := permissionSeeder(tx)
		roles := roleSeeder(tx, permissions)
		userSeeder(tx, roles)

		return nil
	})
}

func permissionSeeder(db *gorm.DB) []entities.Permission {
	var permissions []entities.Permission
	for _, p := range enums.GetAllPermissions() {
		permissions = append(permissions, entities.Permission{
			Name: p.String(),
			Slug: p.String(),
		})
	}

	for _, p := range permissions {
		// Use clause to handle upsert if needed, or just check exists
		var existing entities.Permission
		if err := db.Where("slug = ?", p.Slug).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&p)
			}
		}
	}

	var allPermissions []entities.Permission
	db.Find(&allPermissions)

	fmt.Printf("%d record of permissions seeded\n", len(allPermissions))
	return allPermissions
}

func roleSeeder(db *gorm.DB, permissions []entities.Permission) map[string]entities.Role {
	roles := []entities.Role{
		{Name: "Super Admin"},
		{Name: "User"},
	}

	roleMap := make(map[string]entities.Role)

	for _, r := range roles {
		var existing entities.Role
		if err := db.Where("name = ?", r.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&r)
				existing = r
			}
		}

		if existing.Name == "Super Admin" {
			// Sync all permissions to Super Admin
			db.Model(&existing).Association("Permissions").Replace(permissions)
		}

		roleMap[existing.Name] = existing
	}

	fmt.Printf("%d record of roles seeded\n", len(roles))
	return roleMap
}

func userSeeder(db *gorm.DB, roles map[string]entities.Role) {
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
				db.Create(&u.User)
				existing = u.User
			}
		}

		// Sync Role if RoleName is provided
		if u.RoleName != "" {
			if role, ok := roles[u.RoleName]; ok {
				db.Model(&existing).Association("Roles").Replace([]entities.Role{role})
			}
		}
	}

	fmt.Printf("%d record of users seeded\n", len(users))
}
