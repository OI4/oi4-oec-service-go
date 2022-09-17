package tls

import (
	"crypto/tls"
	"crypto/x509"
)

func NewTLSConfig(ca_pem string, client_cert_pem string, client_cert_key_pem string, tlsVerify bool) (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	client_certificates := make([]tls.Certificate, 0)

	if ca_pem != "" {
		certpool.AppendCertsFromPEM([]byte(ca_pem))
	}

	if client_cert_pem != "" && client_cert_key_pem != "" {
		cert, err := tls.X509KeyPair([]byte(client_cert_pem), []byte(client_cert_key_pem))
		if err != nil {
			return nil, err
		}

		// Just to print out the client certificate..
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return nil, err
		}
		client_certificates = []tls.Certificate{cert}
	}
	// Import client certificate/key pair

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: !tlsVerify,
		// Certificates = list of certs client sends to server.
		Certificates: client_certificates,
	}, nil
}
