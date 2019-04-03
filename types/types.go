package types

type Hero struct {
	Name             string `json:"name"`
	SecretIdentityID int    `json:"secretIdentityID"`
}

type Identity struct {
	ID       int    `json:"id"`
	RealName string `json:"realName"`
}
