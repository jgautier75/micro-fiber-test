package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {

	certPEMBlock, err := os.ReadFile("cert.pem")
	if err != nil {
		fErr := fmt.Errorf("error reading file [%w]", err)
		fmt.Println(fErr)
	}

	var certDERBlock *pem.Block
	certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
	fmt.Printf("Cert block type: %s\n", certDERBlock.Type)

	cert, errDeser := x509.ParseCertificate(certDERBlock.Bytes)
	if errDeser != nil {
		errF := fmt.Errorf("x509_cert_deserialize [%v]", errDeser)
		fmt.Println(errF)
	}

	for _, certNames := range cert.Subject.Names {
		oidTag := getTagForOid(certNames.Type)
		fmt.Printf("Certificate type: [%s], Certificate value [%s]\n", oidTag, certNames.Value)
	}

	for _, issuerName := range cert.Issuer.Names {
		oidTag := getTagForOid(issuerName.Type)
		fmt.Printf("Issuer type: [%s], Issuer value [%s]\n", oidTag, issuerName.Value)
	}

	var buf bytes.Buffer
	fingerprint := sha256.Sum256(cert.Raw)
	for i, f := range fingerprint {
		if i > 0 {
			_, _ = fmt.Fprintf(&buf, ":")
		}
		_, _ = fmt.Fprintf(&buf, "%02X", f)
	}
	fmt.Printf("Fingerprint %s\n", buf.String())
}

func getTagForOid(oid asn1.ObjectIdentifier) string {
	type oidNameMap struct {
		oid  []int
		name string
	}

	oidTags := []oidNameMap{
		{[]int{2, 5, 4, 3}, "CN"},
		{[]int{2, 5, 4, 5}, "SN"},
		{[]int{2, 5, 4, 6}, "C"},
		{[]int{2, 5, 4, 7}, "L"},
		{[]int{2, 5, 4, 8}, "ST"},
		{[]int{2, 5, 4, 10}, "O"},
		{[]int{2, 5, 4, 11}, "OU"},
		{[]int{1, 2, 840, 113549, 1, 9, 1}, "E"}}

	for _, v := range oidTags {
		if oid.Equal(v.oid) {
			return v.name
		}
	}

	return fmt.Sprint(oid)
}
