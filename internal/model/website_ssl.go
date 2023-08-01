package model

import "time"

type WebsiteSSL struct {
	BaseModel
	PrivateKey    string      `gorm:"type:longtext;not null" json:"privateKey"`
	Pem           string      `gorm:"type:longtext;not null" json:"pem"`
	Domain        string      `gorm:"type:varchar(256);not null" json:"domain"`
	DomainPath    string      `gorm:"type:varchar(256);not null" json:"domainPath"`
	CertURL       string      `gorm:"type:varchar(256);not null" json:"certURL"`
	Type          string      `gorm:"type:varchar(64);not null" json:"type"`
	Organization  string      `gorm:"type:varchar(64);not null" json:"organization"`
	AcmeAccountID uint        `gorm:"type:integer;not null" json:"acmeAccountId"`
	ExpireDate    time.Time   `json:"expireDate"`
	StartDate     time.Time   `json:"startDate"`
	AcmeAccount   AcmeAccount `json:"acmeAccount" gorm:"-:migration"`
}

func (w WebsiteSSL) TableName() string {
	return "website_ssl"
}
