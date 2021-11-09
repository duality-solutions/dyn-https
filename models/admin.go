package models

// Admin stores dyn-https admin users and password hash
// hash = SHA256(user + password)
// swagger:parameters models.Admin
type Admin struct {
	UserName     string
	ExpectedHash string
}
