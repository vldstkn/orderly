package api

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type RegisterResponse struct {
	Id          int    `json:"id"`
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Id          int    `json:"id"`
	AccessToken string `json:"access_token"`
}

type GetProfileResponse struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type GetNewTokensResponse struct {
	AccessToken string `json:"access_token"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type ChangeRoleResponse struct {
	Role        string `json:"role"`
	Id          int    `json:"id"`
	AccessToken string `json:"access_token"`
}
