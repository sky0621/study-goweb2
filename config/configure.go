package config

import (
	"flag"
	"log"
)

// "Configure ... "
type Configure struct {
	Domain             string
	Port               string
	SecKey             string
	GoogleClientID     string
	GoogleClientSecret string
}

// "ParseFlag ... "
func ParseFlag() *Configure {
	c := new(Configure)
	var d = flag.String("d", "localhost", "アプリケーションの接続先ドメイン")
	var p = flag.String("p", "80", "アプリケーションの接続先ポート")
	var k = flag.String("k", "secKey", "セキュリティキー")
	var cid = flag.String("c", "googleCid", "GoogleクライアントID")
	var s = flag.String("s", "googleSec", "Googleクライアントシークレット")
	flag.Parse()
	c.Domain = *(d)
	c.Port = *(p)
	c.SecKey = *(k)
	c.GoogleClientID = *(cid)
	c.GoogleClientSecret = *(s)

	log.Println("[flag]domain: " + c.Domain)
	log.Println("[flag]port: " + c.Port)
	log.Println("[flag]secKey: " + c.SecKey)
	log.Println("[flag]googleCid: " + c.GoogleClientID)
	log.Println("[flag]googleSec: " + c.GoogleClientSecret)
	return c
}
