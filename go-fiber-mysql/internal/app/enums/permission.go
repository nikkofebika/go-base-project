package enums

type Permission string

const (
	// User Permissions
	PermissionUserRead        Permission = "user.read"
	PermissionUserCreate      Permission = "user.create"
	PermissionUserUpdate      Permission = "user.update"
	PermissionUserDelete      Permission = "user.delete"
	PermissionUserRestore     Permission = "user.restore"
	PermissionUserForceDelete Permission = "user.force_delete"
	PermissionUserSyncRoles   Permission = "user.sync_roles"

	// Role Permissions
	PermissionRoleRead            Permission = "role.read"
	PermissionRoleCreate          Permission = "role.create"
	PermissionRoleUpdate          Permission = "role.update"
	PermissionRoleDelete          Permission = "role.delete"
	PermissionRoleRestore         Permission = "role.restore"
	PermissionRoleForceDelete     Permission = "role.force_delete"
	PermissionRoleSyncPermissions Permission = "role.sync_permissions"
)

func (p Permission) String() string {
	return string(p)
}

func GetAllPermissions() []Permission {
	return []Permission{
		PermissionUserRead,
		PermissionUserCreate,
		PermissionUserUpdate,
		PermissionUserDelete,
		PermissionUserRestore,
		PermissionUserForceDelete,
		PermissionUserSyncRoles,

		PermissionRoleRead,
		PermissionRoleCreate,
		PermissionRoleUpdate,
		PermissionRoleDelete,
		PermissionRoleRestore,
		PermissionRoleForceDelete,
		PermissionRoleSyncPermissions,
	}
}
