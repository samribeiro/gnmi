package credentials

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"

	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	ca             = flag.String("ca", "", "CA certificate file.")
	cert           = flag.String("cert", "", "Certificate file.")
	key            = flag.String("key", "", "Private key file.")
	authorizedUser = userCredentials{}
	usernameKey    = "username"
	passwordKey    = "password"
)

func init() {
	flag.StringVar(&authorizedUser.username, "username", "", "If specified, uses username/password credentials.")
	flag.StringVar(&authorizedUser.password, "password", "", "The password matching the provided username.")
}

type userCredentials struct {
	username string
	password string
}

func (a *userCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		usernameKey: a.username,
		passwordKey: a.password,
	}, nil
}

func (a *userCredentials) RequireTransportSecurity() bool {
	return true
}

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

func ClientCredentials(server string) []grpc.DialOption {
	opts := []grpc.DialOption{}

	certificates, certPool := LoadCertificates()
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		ServerName:   server, // This is required and must match the certificate CN.
		Certificates: certificates,
		RootCAs:      certPool,
	})))

	if authorizedUser.username != "" {
		return append(opts, grpc.WithPerRPCCredentials(&authorizedUser))
	}
	return opts
}

func ServerCredentials() []grpc.ServerOption {
	certificates, certPool := LoadCertificates()

	return []grpc.ServerOption{grpc.Creds(credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: certificates,
		ClientCAs:    certPool,
	}))}
}

func AuthorizeUser(ctx context.Context) (string, bool) {
	authorize := false
	if authorizedUser.username == "" {
		authorize = true
	}
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "no Metadata found", authorize
	}
	user, ok := headers[usernameKey]
	if !ok && len(user) > 0 {
		return "no username in Metadata", authorize
	}
	pass, ok := headers[passwordKey]
	if !ok && len(pass) > 0 {
		return "found username \"%s\" but no password in Metadata", authorize
	}
	if authorize || pass[0] == authorizedUser.password && user[0] == authorizedUser.username {
		return fmt.Sprintf("authorized with \"%s:%s\"", user[0], pass[0]), true
	}
	return fmt.Sprintf("not authorized with \"%s:%s\"", user[0], pass[0]), false
}
