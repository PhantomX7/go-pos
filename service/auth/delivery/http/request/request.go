package request

type SignInRequest struct {
	Username string `form:"username" binding:"required,exist=users.username"`
	Password string `form:"password" binding:"required"`
}

type SignUpRequest struct {
	Username   string `form:"username" binding:"required,unique=users.username"`
	Password   string `form:"password" binding:"required"`
	Passphrase string `form:"passphrase" binding:"required"`
	RoleID     uint64 `form:"role_id" binding:"required,exist=roles.id"`
}
