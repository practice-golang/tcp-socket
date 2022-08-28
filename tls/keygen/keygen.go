package main // import "keygen"

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"time"
)

func getPrivateKey() (keyString string, publickey crypto.Signer, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", nil, err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", nil, err
	}

	keyBlock := pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}
	keyString = string(pem.EncodeToMemory(&keyBlock))

	return keyString, key, nil
}

func getPublicKey(privateString string, privKey crypto.Signer, template *x509.Certificate) (publicKeyString string, cert string, err error) {
	privateBlock, rest := pem.Decode([]byte(privateString))
	if len(rest) > 0 {
		log.Fatal(len(rest))
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		log.Fatal("ParsePKCS8PrivateKey:", err)
	}

	if reflect.TypeOf(privateKey).String() != "*rsa.PrivateKey" {
		return "", "", fmt.Errorf("pkey is not *rsa.PrivateKey")
	}

	publicKey := privateKey.(*rsa.PrivateKey).Public()
	if reflect.TypeOf(publicKey).String() != "*rsa.PublicKey" {
		return "", "", fmt.Errorf("not rsa")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyBlock := pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}
	publicKeyString = string(pem.EncodeToMemory(&publicKeyBlock))

	// dkim = "v=DKIM1;k=rsa;p=" + base64.StdEncoding.EncodeToString(publicKeyBytes)
	// dkim = "k=rsa;p=" + base64.StdEncoding.EncodeToString(publicKeyBytes)

	parent := template
	certByte, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privKey)
	if err != nil {
		return "", "", err
	}

	// cert = string(certByte)
	cert = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certByte}))

	return publicKeyString, cert, err
}

func writeToFile(data, filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func getCertTemplate() *x509.Certificate {
	// http://golang.org/pkg/crypto/x509/#KeyUsage
	template := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Subject:               pkix.Name{Country: []string{"My home"}, Organization: []string{"My room"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	return template
}

func main() {
	template := getCertTemplate()

	privString, privKey, err := getPrivateKey()
	if err != nil {
		log.Fatal("generatePrivateKEY:", err)
	}

	pubString, cert, err := getPublicKey(privString, privKey, template)
	if err != nil {
		log.Fatal("getDKIM:", err)
	}

	err = writeToFile(privString, "key.pem")
	if err != nil {
		log.Fatal("writeToFile: key.pem", err)
	}

	err = writeToFile(pubString, "key.pub")
	if err != nil {
		log.Fatal("writeToFile: key.pub", err)
	}

	err = writeToFile(cert, "key.crt")
	if err != nil {
		log.Fatal("writeToFile: key.crt", err)
	}

	fmt.Println(privString)
	fmt.Println(pubString)
	fmt.Println(cert)
}
