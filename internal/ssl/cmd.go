package ssl

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hzchiyan/cy/internal/model"
	"gorm.io/gorm"
	"log"
)

func WebsiteSSLAccount(db *gorm.DB, email string, privateKeys ...string) (aa model.AcmeAccount, err error) {
	privateKey := ""
	if len(privateKeys) > 0 {
		privateKey = privateKeys[0]
	}
	_ = db.Model(model.AcmeAccount{}).Where(model.AcmeAccount{Email: email}).First(&aa).Error
	if aa.ID > 0 {
		log.Printf("在db中查到数据 %v", aa)
		return aa, nil
	}
	log.Printf("正在生成用户数据")
	client, err := NewAcmeClient(email, privateKey)
	if err != nil {

		return aa, err
	}
	aa.Email = email
	aa.URL = client.User.Registration.URI
	aa.PrivateKey = string(GetPrivateKey(client.User.GetPrivateKey()))
	if err = db.Save(&aa).Error; err != nil {
		log.Printf("在保存用户数据数据 %v", aa)
		return aa, err
	}
	return aa, nil
}

func WebsiteSSL(db *gorm.DB, email, domain, domainPath string) (websiteSSL model.WebsiteSSL, err error) {
	account, err := WebsiteSSLAccount(db, email)
	if err != nil {
		log.Printf("查下注册用户信息失败")
		return websiteSSL, err
	}
	websiteSSL.AcmeAccountID = account.ID
	websiteSSL.Domain = domain
	websiteSSL.DomainPath = domainPath
	_ = db.Model(model.WebsiteSSL{}).Where(websiteSSL).First(&websiteSSL).Error
	if websiteSSL.ID > 0 {
		log.Printf("在db中查到数据 %v", websiteSSL)
		if websiteSSL.StartDate.Before(websiteSSL.ExpireDate) {
			return websiteSSL, nil
		}
	}
	log.Printf("正在重新生成")
	client, err := NewPrivateKeyClient(account.Email, account.PrivateKey)
	if err != nil {
		return
	}
	if err := client.UseHTTP(domainPath); err != nil {
		log.Printf("client.UseHTTP err=%v", err)
		return websiteSSL, fmt.Errorf("client.UseHTTP err=%v", err)
	}
	resource, err := client.ObtainSSL([]string{domain})
	if err != nil {
		log.Printf("client.ObtainSSL err=%v", err)
		return websiteSSL, err
	}
	websiteSSL.PrivateKey = string(resource.PrivateKey)
	websiteSSL.Pem = string(resource.Certificate)
	websiteSSL.CertURL = resource.CertURL
	certBlock, _ := pem.Decode(resource.Certificate)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		log.Printf("x509.ParseCertificate err=%v", err)
		return websiteSSL, err
	}
	websiteSSL.ExpireDate = cert.NotAfter
	websiteSSL.StartDate = cert.NotBefore
	websiteSSL.Type = cert.Issuer.CommonName
	websiteSSL.Organization = cert.Issuer.Organization[0]
	if err = db.Save(&websiteSSL).Error; err != nil {
		log.Printf("保存站点信息失败 err=%v", websiteSSL)
		return websiteSSL, err
	}
	return websiteSSL, nil
}
