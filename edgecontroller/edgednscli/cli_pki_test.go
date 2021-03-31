// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cli_test

// Starting client to Edgednssvr
import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"google.golang.org/grpc/credentials"
)

var (
	// CAKey generated with:
	// openssl genpkey -algorithm RSA -out cakey.pem

	// this key is here just to allow to regenerate the rest
	// 	CAKey = []byte(`-----BEGIN PRIVATE KEY-----
	// MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAOtKSk/bKAOY2wKI
	// i936z2JjWRPeIOePYa278lkPin3oXr5SmiZTVEbmV6hoL88B9UXP9291N5PsKAfF
	// RPzhcW8w7MJgmkmoWch+iqfeme3tbDRnVcDzOs1v3BnQDEEffoTccnNuX9Ji0Xq1
	// oFH5ys4jy+fEVh2bBSe8Nm/zMMdlAgMBAAECgYEAxWrT5yTu8fWy7PE8ycfDcb4B
	// 45Wg7ilzv53OvcM4ZGQv+mi8RiPstoBWVEr2si3maekQp0gmlcqz27E9LSDLVvBD
	// eTRb48JqOxenTHZVPj9AO1n9CPlhV6XiUqdbu328LTi4VCgCYLLkQlJ1/zQAotbi
	// mggEdgwN8V15LLpblgECQQD5Q74E8U++9TlwAebkfgr+xUUAa011GifuMUgb5dVi
	// NSpPwA821s/OYG5ksG70pO7VDUxpZ/QTWT6LlJdjXW/FAkEA8aXhs+ymHpZZAZlm
	// zOxlo7RjDe/bPp1pJ8mBkHo1K106fx0vm4Zcx7REzarYrOqfJ/wmK/2b7Ixbs2xL
	// 907TIQI/Vqg2Z4F9N+YJEeeXXdhz+h6a3opvY0a0mxqG9yvwAzR4gUFECMylchoU
	// Ru0Ppf2xCmCVMA/zysY7+EqwVtslAkEAs6L8ETQYHd7olb7C5A3MVRERHdgz5O1B
	// sbkjnxKqGe15Qeuha31L1Y2WJEogPxM/c4sHmWYzlnqU5OHZ6yPm4QJBAPKsMron
	// 7YN7OzQze3Uu0UQoOEfTRTrybgdodz/YhbGVfj3KIK4cS9v2YVusXcAEzw+hgfWZ
	// najMYyTaQOUOoX8=
	// -----END PRIVATE KEY-----`)

	// CACert generated with:
	// openssl req -new -x509 -key cakey.pem -days 36500 \
	// -subj "/C=TT/ST=TestCA/L=TestCA/O=TestCA/OU=TestCA/CN=TestCA" \
	// -out cacert.pem
	CACert = []byte(`-----BEGIN CERTIFICATE-----
MIIClDCCAf2gAwIBAgIJAPrpCwima9g8MA0GCSqGSIb3DQEBCwUAMGIxCzAJBgNV
BAYTAlRUMQ8wDQYDVQQIDAZUZXN0Q0ExDzANBgNVBAcMBlRlc3RDQTEPMA0GA1UE
CgwGVGVzdENBMQ8wDQYDVQQLDAZUZXN0Q0ExDzANBgNVBAMMBlRlc3RDQTAgFw0x
OTA3MDQwNjU1NTNaGA8yMTE5MDYxMDA2NTU1M1owYjELMAkGA1UEBhMCVFQxDzAN
BgNVBAgMBlRlc3RDQTEPMA0GA1UEBwwGVGVzdENBMQ8wDQYDVQQKDAZUZXN0Q0Ex
DzANBgNVBAsMBlRlc3RDQTEPMA0GA1UEAwwGVGVzdENBMIGfMA0GCSqGSIb3DQEB
AQUAA4GNADCBiQKBgQDrSkpP2ygDmNsCiIvd+s9iY1kT3iDnj2Gtu/JZD4p96F6+
UpomU1RG5leoaC/PAfVFz/dvdTeT7CgHxUT84XFvMOzCYJpJqFnIfoqn3pnt7Ww0
Z1XA8zrNb9wZ0AxBH36E3HJzbl/SYtF6taBR+crOI8vnxFYdmwUnvDZv8zDHZQID
AQABo1AwTjAdBgNVHQ4EFgQU8zdc8gyws69MSS82zfSa1Smfo8swHwYDVR0jBBgw
FoAU8zdc8gyws69MSS82zfSa1Smfo8swDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOBgQBaZ/TwdeE2Jh1bY5JP0xg8Js6z2db5y9pcV6ywCppbCqm+55dsgjVV
s22+o2g0qPu4ukilpOq51JCN2XvTKGqbUTGs1C9oQQDtNv4rBwYwunv6cRfJ6qdn
AGMWVUY2WBQenzG86L7dbiaqziq7jaTPyjLdfRmdf89B1l+VAg7QPg==
-----END CERTIFICATE-----`)
	// ServerKey generated with:
	// openssl genpkey -algorithm RSA -out serverkey.pem
	ServerKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAPNp+vH8grjgJhSE
HLN2leLcHEmepQWIYwL6p4E7jBP4xR84i0V9ghhND262XhrPwWmOF8neYcvQwFuv
OtoMgVr1cu4LXlJ8D4VD/iCLaAtt6acl8iDYKd8rpSRR5fS4QmRem/4ABVllLrq5
3L5CDm77Q65UdPaG6RU0LLfm+p15AgMBAAECgYBpW/Mp8QFoImV5NeHycB7ajimU
WB3XsjwhGvuL/YlZT7GLHP9zearo99n5MuiytEBkfwOe9ingfZ/1T3iqzwTh0QgX
67ivTENOfixYIAFxnXICxqMedHbSf4g7I84MvYmuij+AwfmKynxHihAyePgUkgpT
EPEYShvONJiaJv8EzQJBAPo5SO+Q1xv3FCkZKGXCcnemBZc5PcIXLxAadylhDPba
G4PcwvCkcJQAE+mYUoeqgWjWHPpI0IrcDKrsqfMu+k8CQQD5CHPTvGqSdMY/X4nR
irCNnl8DdLuP9yLg/iPElFB4qA9+tHt9vxCdf3UYFlGHO+0eonpWUUhLTcgkZwG+
5qG3AkBPLMsgSqNWDPNLIdO/hyxzIYxJUmpxPKa+oLEt3bCdd1wXeEYSoOWiXvfu
BS2wEcqK4p0esOWT2pcDiokoEqhPAkEAnM9Qt+v1o2D6kF8lVhehkps9stxnWoDA
qrH8bcUCYudYUm1tMUIFQVdHsTIkhiYa6mZe6EloX1qMAzgFZmKtQwJAGpoJfnz/
xwLN4/bq14PdCfZuUBXBmIz/QNfDt4lWZ5KUf8GBQtqRo3nwbKODqiP8URicF5S3
u7w2aa5w9B9vwg==
-----END PRIVATE KEY-----`)
	// ServerCert generated with:
	// openssl req -new -key serverkey.pem -out server.csr \
	// -subj "/C=TT/ST=Server/L=Server/O=Server/OU=Server/CN=localhost"
	// openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost") -in server.csr \
	// -CA cacert.pem -CAkey cakey.pem -days 36500 -out servercert.pem -CAcreateserial
	ServerCert = []byte(`-----BEGIN CERTIFICATE-----
MIICajCCAdOgAwIBAgIUVy/z7c0lDyiB45G+ut0s/aytvAcwDQYJKoZIhvcNAQEL
BQAwYjELMAkGA1UEBhMCVFQxDzANBgNVBAgMBlRlc3RDQTEPMA0GA1UEBwwGVGVz
dENBMQ8wDQYDVQQKDAZUZXN0Q0ExDzANBgNVBAsMBlRlc3RDQTEPMA0GA1UEAwwG
VGVzdENBMCAXDTIxMDIxMDEyMDQzNFoYDzIxMjEwMTE3MTIwNDM0WjBlMQswCQYD
VQQGEwJUVDEPMA0GA1UECAwGU2VydmVyMQ8wDQYDVQQHDAZTZXJ2ZXIxDzANBgNV
BAoMBlNlcnZlcjEPMA0GA1UECwwGU2VydmVyMRIwEAYDVQQDDAlsb2NhbGhvc3Qw
gZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAPNp+vH8grjgJhSEHLN2leLcHEme
pQWIYwL6p4E7jBP4xR84i0V9ghhND262XhrPwWmOF8neYcvQwFuvOtoMgVr1cu4L
XlJ8D4VD/iCLaAtt6acl8iDYKd8rpSRR5fS4QmRem/4ABVllLrq53L5CDm77Q65U
dPaG6RU0LLfm+p15AgMBAAGjGDAWMBQGA1UdEQQNMAuCCWxvY2FsaG9zdDANBgkq
hkiG9w0BAQsFAAOBgQBG5xF3jdfcW8EkdeLTLa2wy5hzE+35p2pwOFUP7zkjkCaq
2QcbvhkcoAR9KomH7qa2JPino1+2O5R0ATzv4CwzKpQvRAm+jVroxuTP3TI4yy4j
tqIWTyZrFqr1w/5yCh6lE5Panv/FEezv9WguBnFG04Vj0c/toUTqxCKKjH613w==
-----END CERTIFICATE-----`)
	// ClientKey generated with:
	// openssl genpkey -algorithm RSA -out clientkey.pem
	ClientKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMA9zjnHCs5dyoVp
IwQi60kOns5FywrzT/sPEq6Fea19tGmcRN0KVSqt0Oi4/Uaw6+utAlpO8T2KVgN3
Ds1TxbtiTuy5BuQnQwlZoxwFfPQCQoOqF6e77IM4OHCkzR2d2RgP9liPejKqy9KY
ggjVJSpyI5q+pg3UsB39y9IEPWUvAgMBAAECgYA9AeA2WcJGJ2Gj4O+p2mMMtkUT
dPawrPKO3AOzuswYkLH2ag84jpTkjHClyj8mdfMIwcjjW/YW75XIV7os0w8GkQvL
0KyUx0GClEzPcLK3FN5T5sAb/BKk3MozXq6uMUlcuzGKtJQZo550xStPFfBUzdx/
++C0j9Vx0R1jTm7WsQJBAPOBfid/oHwCTPhHcXz7FWZTI2tzogr5D8zpIGRaAW+s
x6xGE19nGVMOC3YRRvPLFqTkWwnMkSDD661R1xxb57UCQQDKGvFOvEw2nERFyzAQ
40KWQXErHmNVvAZYSrU/XahBepknE4fAsajE63z9j+VKNL9o45lPdxWlcoPqUpYl
xJ/TAkEA0s00iD8tkI1xUrKOpdK7mOS6Ugg9rOAmdI8RMq7UdXgPN61Gkh2vx73t
4B1e3ijfhooF4frNXJ8fhODHp05MzQJASv+0jINTxPA4Za4BKEMtcpTWaincsvQT
ZrBSkjEvPx0/uECHH/rNpkW9HFtJuz/XyBStz89hOQRF7662CU3uAQJAFILl2yD1
JEYobtnN2vxSQJ5cJLKMo+AmdruHUshnbQFn/FWegEKYjbNm0c4Uh7VlWaGl1NK2
rVvUqyaq0a6JJA==
-----END PRIVATE KEY-----`)
	// ClientCert generated with:
	// 	openssl req -new -key clientkey.pem -out client.csr \
	// 	-subj "/C=TT/ST=Client/L=Client/O=Client/OU=Client/CN=client"
	//    openssl x509 -req -in client.csr -CA cacert.pem -CAkey cakey.pem \
	// 	-days 36500 -out clientcert.pem -CAcreateserial
	ClientCert = []byte(`-----BEGIN CERTIFICATE-----
MIICPTCCAaYCCQDx/yry8DvrADANBgkqhkiG9w0BAQsFADBiMQswCQYDVQQGEwJU
VDEPMA0GA1UECAwGVGVzdENBMQ8wDQYDVQQHDAZUZXN0Q0ExDzANBgNVBAoMBlRl
c3RDQTEPMA0GA1UECwwGVGVzdENBMQ8wDQYDVQQDDAZUZXN0Q0EwIBcNMTkwNzA0
MDcwNDA1WhgPMjExOTA2MTAwNzA0MDVaMGIxCzAJBgNVBAYTAlRUMQ8wDQYDVQQI
DAZDbGllbnQxDzANBgNVBAcMBkNsaWVudDEPMA0GA1UECgwGQ2xpZW50MQ8wDQYD
VQQLDAZDbGllbnQxDzANBgNVBAMMBmNsaWVudDCBnzANBgkqhkiG9w0BAQEFAAOB
jQAwgYkCgYEAwD3OOccKzl3KhWkjBCLrSQ6ezkXLCvNP+w8SroV5rX20aZxE3QpV
Kq3Q6Lj9RrDr660CWk7xPYpWA3cOzVPFu2JO7LkG5CdDCVmjHAV89AJCg6oXp7vs
gzg4cKTNHZ3ZGA/2WI96MqrL0piCCNUlKnIjmr6mDdSwHf3L0gQ9ZS8CAwEAATAN
BgkqhkiG9w0BAQsFAAOBgQBVGopOMrm6FzxRR9+NDTApFDpCv9lOV7XVsY3WTuIE
3O2RQt4epjzjzgCBmPf/+PV5g38CXgK/+urHvwei0DjJeBtCk4kTppfy/KhOd7qm
v+f9sIm0cc/ErZLqA0nGlb76TfjLrcabOYqXgAtRLtSUpxToutdvwZlE0gKHCN8L
xQ==
-----END CERTIFICATE-----`)
)

// prepareTestCredentials performs fake grpc credential preparation by saving
// hardcoded keys and certificates.
// This should be used only for testing
func prepareTestCredentials(certsDir string) error {
	if err := os.MkdirAll(certsDir, 0700); err != nil {
		return fmt.Errorf("Failed to create %s: %v", certsDir, err)
	}

	if err := ioutil.WriteFile(filepath.Join(certsDir, "key.pem"),
		ServerKey, 0600); err != nil {
		return fmt.Errorf("Failed to save key: %v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(certsDir, "cert.pem"),
		ServerCert, 0600); err != nil {
		return fmt.Errorf("Failed to save cert: %v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(certsDir, "cacerts.pem"),
		CACert, 0600); err != nil {
		return fmt.Errorf("Failed to save cacerts: %v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(certsDir, "root.pem"),
		CACert, 0600); err != nil {
		return fmt.Errorf("Failed to save root ca: %v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(certsDir, "c_key.pem"),
		ClientKey, 0600); err != nil {
		return fmt.Errorf("Failed to save client key: %v", err)
	}

	if err := ioutil.WriteFile(filepath.Join(certsDir, "c_cert.pem"),
		ClientCert, 0600); err != nil {
		return fmt.Errorf("Failed to save client cert: %v", err)
	}

	return nil
}

// readPKI reads pki credentials just like regular server does.
// this should be used only for testing
func readTestPKICredentials(crtPath, keyPath,
	caPath string) (*credentials.TransportCredentials, error) {

	srvCert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed load server key pair: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("Failed read ca certificates: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Printf("Failed appends CA certs from %s", caPath)
		return nil, fmt.Errorf("Failed appends CA certs from %s", caPath)
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{srvCert},
		ClientCAs:    certPool,
	})

	return &creds, nil
}
