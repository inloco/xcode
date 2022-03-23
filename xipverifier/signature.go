package main

// #include "xipverifier.h"
import "C"

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"unsafe"

	cms "github.com/github/smimesign/ietf-cms"
)

type XARSignature struct {
	c C.xar_signature_t
}

func (xs *XARSignature) C() C.xar_signature_t {
	return xs.c
}

func (xs *XARSignature) IsInitialized() bool {
	return xs.c != nil
}

func (xs *XARSignature) GetType() (string, error) {
	if !xs.IsInitialized() {
		return "", errors.New("!xs.IsInitialized()")
	}

	c := C.xar_signature_type(xs.c)
	t := C.GoString(c)

	return t, nil
}

func (xs *XARSignature) GetX509Certificate(index int) (*x509.Certificate, error) {
	if !xs.IsInitialized() {
		return nil, errors.New("!xs.IsInitialized()")
	}

	var certData *C.uint8_t
	var certSize C.uint32_t
	if err := C.xar_signature_get_x509certificate_data(xs.c, C.int32_t(index), &certData, &certSize); err != 0 {
		return nil, GetError()
	}

	cert := C.GoBytes(unsafe.Pointer(certData), C.int(certSize))
	return x509.ParseCertificate(cert)
}

func (xs *XARSignature) GetX509Certificates() ([]*x509.Certificate, error) {
	if !xs.IsInitialized() {
		return nil, errors.New("!xs.IsInitialized()")
	}

	count := C.xar_signature_get_x509certificate_count(xs.c)
	if count < 0 {
		return nil, GetError()
	}

	var certificates []*x509.Certificate
	for i := 0; i < int(count); i++ {
		certificate, err := xs.GetX509Certificate(i)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certificate)
	}

	return certificates, nil
}

func (xs *XARSignature) CopySignedData() ([]byte, []byte, error) {
	if !xs.IsInitialized() {
		return nil, nil, errors.New("!xs.IsInitialized()")
	}

	var rawData *C.uint8_t
	var rawSize C.uint32_t
	var sigData *C.uint8_t
	var sigSize C.uint32_t
	if err := C.xar_signature_copy_signed_data(xs.c, &rawData, &rawSize, &sigData, &sigSize, nil); err != 0 {
		return nil, nil, GetError()
	}
	defer C.free(unsafe.Pointer(rawData))
	defer C.free(unsafe.Pointer(sigData))

	raw := C.GoBytes(unsafe.Pointer(rawData), C.int(rawSize))
	sig := C.GoBytes(unsafe.Pointer(sigData), C.int(sigSize))

	return raw, sig, nil
}

func (xs *XARSignature) Verify() ([]byte, error) {
	if !xs.IsInitialized() {
		return nil, errors.New("!xs.IsInitialized()")
	}

	t, err := xs.GetType()
	if err != nil {
		return nil, err
	}

	switch t {
	case "CMS":
		return xs.verifyCMS()

	case "RSA":
		return xs.verifyRSA()

	default:
		return nil, fmt.Errorf("unknown xar signature type: %s", t)
	}
}

func (xs *XARSignature) verifyCMS() ([]byte, error) {
	raw, sig, err := xs.CopySignedData()
	if err != nil {
		return nil, err
	}

	sd, err := cms.ParseSignedData(sig)
	if err != nil {
		return nil, err
	}

	opts := x509.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageAny,
		},
	}
	if _, err := sd.VerifyDetached(raw, opts); err != nil {
		return nil, err
	}

	return raw, nil
}

func (xs *XARSignature) verifyRSA() ([]byte, error) {
	leaf, err := xs.verifyLeaf()
	if err != nil {
		return nil, err
	}

	if alg := leaf.PublicKeyAlgorithm; alg != x509.RSA {
		return nil, fmt.Errorf(`%s != x509.RSA`, alg)
	}

	pub, ok := leaf.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("!ok")
	}

	raw, sig, err := xs.CopySignedData()
	if err != nil {
		return nil, err
	}

	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA1, raw, sig); err != nil {
		return nil, err
	}

	return raw, nil
}

func (xs *XARSignature) verifyLeaf() (*x509.Certificate, error) {
	certs, err := xs.GetX509Certificates()
	if err != nil {
		return nil, err
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("len(certs) == 0")
	}

	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageAny,
		},
	}
	for _, cert := range certs {
		opts.Intermediates.AddCert(cert)
	}

	chains, err := certs[0].Verify(opts)
	if err != nil {
		return nil, err
	}

	return chains[0][0], nil
}
