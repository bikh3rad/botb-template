package dto

import (
	"time"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/entity"
)

const rfc3339 = "2006-01-02T15:04:05Z07:00"

// LoginReq is the login request body.
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshReq carries a refresh token (refresh and logout).
type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

// AdminResp is the API representation of an admin account. The password hash
// deliberately has no field here — it never leaves the service.
type AdminResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	LastLoginAt string `json:"last_login_at,omitempty"`
}

// ToAdminResp maps an account entity to its API shape.
func ToAdminResp(a entity.AdminAccount) AdminResp {
	resp := AdminResp{
		ID:        a.ID.String(),
		Name:      a.Name,
		Email:     a.Email,
		Role:      string(a.Role),
		IsActive:  a.IsActive,
		CreatedAt: formatTime(a.CreatedAt),
	}

	if a.LastLoginAt != nil {
		resp.LastLoginAt = formatTime(*a.LastLoginAt)
	}

	return resp
}

// TokenResp is the login/refresh response: a short-lived access JWT plus a
// one-time-use refresh token (rotated on every refresh).
type TokenResp struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Admin        AdminResp `json:"admin"`
}

// ToTokenResp maps a login result to the API shape.
func ToTokenResp(r biz.LoginResult) TokenResp {
	return TokenResp{
		AccessToken:  r.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    r.ExpiresIn,
		RefreshToken: r.RefreshToken,
		Admin:        ToAdminResp(r.Admin),
	}
}

// AccountCreateReq is the superadmin account-creation body.
type AccountCreateReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// AccountUpdateReq is a partial account edit; omitted fields are unchanged.
type AccountUpdateReq struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

// ToUpdateInput maps the request to the biz input.
func (r AccountUpdateReq) ToUpdateInput() biz.UpdateAccountInput {
	input := biz.UpdateAccountInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		IsActive: r.IsActive,
	}

	if r.Role != nil {
		role := entity.Role(*r.Role)
		input.Role = &role
	}

	return input
}

// AccountListResp is the account list envelope.
type AccountListResp struct {
	Count    int         `json:"count"`
	Accounts []AdminResp `json:"accounts"`
}

// ToAccountListResp maps accounts to the list envelope.
func ToAccountListResp(accounts []entity.AdminAccount) AccountListResp {
	items := make([]AdminResp, 0, len(accounts))
	for i := range accounts {
		items = append(items, ToAdminResp(accounts[i]))
	}

	return AccountListResp{Count: len(items), Accounts: items}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(rfc3339)
}
