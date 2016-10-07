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
	var sk = flag.String("sk", "secKey", "セキュリティキー")
	var gc = flag.String("gc", "googleCid", "GoogleクライアントID")
	var gs = flag.String("gs", "googleSec", "Googleクライアントシークレット")
	// ↑で、いったんコマンドライン引数をアドレスとして取得（定義）しないと、↓のパースが機能しない・・・。
	flag.Parse()
	c.Domain = *(d)
	c.Port = ":" + *(p)
	c.SecKey = *(sk)
	c.GoogleClientID = *(gc)
	c.GoogleClientSecret = *(gs)

	log.Println("[flag]domain: " + c.Domain)
	log.Println("[flag]port: " + c.Port)
	log.Println("[flag]secKey: " + c.SecKey)
	log.Println("[flag]googleCid: " + c.GoogleClientID)
	log.Println("[flag]googleSec: " + c.GoogleClientSecret)
	return c
}
