package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"os"
	"path"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}

func (u *MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func GenerateCert(storageDir, domain string) (*tls.Certificate, error) {
	// detect if cert exists
	dirPath := path.Join(storageDir, domain)
	certPath := path.Join(dirPath, fmt.Sprintf("%s.crt", domain))
	keyPath := path.Join(dirPath, fmt.Sprintf("%s.key", domain))
	_, err := os.Stat(dirPath)
	// has domain
	if err == nil {
		cert, err := os.ReadFile(certPath)
		if err != nil {
			return nil, err
		}
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
		tlsCert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		return &tlsCert, nil
	}

	// generate
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	myUser := MyUser{
		Email: "akazwz@outlook.com",
		Key:   privateKey,
	}

	config := lego.NewConfig(&myUser)
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	dnsProvider, err := dns.NewDNSChallengeProviderByName("alidns")
	if err != nil {
		return nil, err
	}
	err = client.Challenge.SetDNS01Provider(dnsProvider)
	if err != nil {
		return nil, err
	}
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	myUser.Registration = reg
	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}
	// save cert
	cert := certificates.Certificate
	key := certificates.PrivateKey
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(certPath, cert, 0644)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(keyPath, key, 0644)
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return &tlsCert, nil
}
