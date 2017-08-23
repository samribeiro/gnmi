package helper

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"

	"google.golang.org/grpc/credentials"

	log "github.com/golang/glog"
)

var (
	ca   = flag.String("ca", "", "CA certificate file.")
	cert = flag.String("cert", "", "Certificate file.")
	key  = flag.String("key", "", "Private key file.")
)

func LoadCertificates() ([]tls.Certificate, *x509.CertPool) {
	if *ca == "" || *cert == "" || *key == "" {
		log.Exit("-ca -cert and -key must be set with file locations")
	}

	certificate, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Exitf("could not load client key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	caFile, err := ioutil.ReadFile(*ca)
	if err != nil {
		log.Exitf("could not read CA certificate: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(caFile); !ok {
		log.Exit("failed to append CA certificate")
	}

	return []tls.Certificate{certificate}, certPool
}

func ClientCertificates(server string) credentials.TransportCredentials {
	certificates, certPool := LoadCertificates()
	return credentials.NewTLS(&tls.Config{
		ServerName:   server, // This is required and must match the certificate CN.
		Certificates: certificates,
		RootCAs:      certPool,
	})
}

func ServerCertificates() credentials.TransportCredentials {
	certificates, certPool := LoadCertificates()
	return credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: certificates,
		ClientCAs:    certPool,
	})
}
