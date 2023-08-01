package model

type AcmeAccount struct {
	BaseModel
	Email      string `gorm:"type:varchar(256);not null" json:"email"`
	URL        string `gorm:"type:varchar(256);not null" json:"url"`
	PrivateKey string `gorm:"type:longtext;not null" json:"-"`
}

func (w AcmeAccount) TableName() string {
	return "acme_account"
}
