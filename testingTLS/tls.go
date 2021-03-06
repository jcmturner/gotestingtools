package testingTLS

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"testing"
	"time"
)

func GenerateSelfSignedTLSKeyPairFiles(t *testing.T) (string, string, []byte, *rsa.PrivateKey) {
	derBytes, priv := GenerateSelfSignedTLSKeyPairData(t)
	certOut, _ := ioutil.TempFile(os.TempDir(), "testCert")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	keyOut, _ := ioutil.TempFile(os.TempDir(), "testKey")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	return certOut.Name(), keyOut.Name(), derBytes, priv
}

func GenerateSelfSignedTLSKeyPairData(t *testing.T) ([]byte, *rsa.PrivateKey) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 2 * 365 * 24)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}
	template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"))
	template.DNSNames = append(template.DNSNames, "testhost.example.com")
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Errorf("Error creating certifcate for testing: %v", err)
	}
	return derBytes, priv
}

// Cannot write a cert pool to a file as the certificates are not exported from the struct in any way
//func WriteCertPoolToFile(t *testing.T, cp x509.CertPool) (*os.File) {
//	certOut, _ := ioutil.TempFile(os.TempDir(), "testCert")
//	for _, c := range cp.Subjects() {
//		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: c})
//	}
//	certOut.Close()
//	return certOut
//}

func WriteCertToFile(t *testing.T, c *x509.Certificate) *os.File {
	certOut, _ := ioutil.TempFile(os.TempDir(), "testCert")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
	certOut.Close()
	return certOut
}
