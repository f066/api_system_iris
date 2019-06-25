package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/fullsailor/pkcs7"
)

/*
暂未解决 Go语言smime签名
 */

var nl = []byte{'\n'}

type Smime struct {
	Msg		[]byte
	Bound 	string
	Pkcs7 	[]byte
	Cert	*x509.Certificate
	Pkey 	crypto.PrivateKey
}

type LineBreaker struct {
	line 	[]byte
	used	int
	out 	*bytes.Buffer
}

func (this *LineBreaker) SetWriter(out *bytes.Buffer) {
	this.out = out
}

func (this *LineBreaker) Handle(b []byte,leng int) int {
	if this.used+len(b) < leng {
		copy(this.line[this.used:], b)
		this.used += len(b)
		this.out.Write(b)
		return len(b)
	}

	this.out.Write(this.line[0:this.used])

	excess := leng - this.used
	this.used = 0

	this.out.Write(b[0:excess])

	this.out.Write(nl)

	return this.Handle(b[excess:],leng)

}

func (this *Smime) New(cert *x509.Certificate, pkey crypto.PrivateKey) *Smime {
	this.Cert, this.Pkey = cert, pkey
	return this
}

func (this *Smime) SignMsg(msg []byte) {
	this.Msg = msg
	this.Sign()
}

func (this *Smime) Sign() {
	sd ,_ := pkcs7.NewSignedData(this.Msg)
	sd.AddSigner(this.Cert,this.Pkey,pkcs7.SignerInfoConfig{})
	sd.Detach()
	res,err := sd.Finish()
	if err == nil {
		this.Pkcs7 = res
	}else {
		fmt.Printf("Cannot finish signing data: %s", err)
	}
}

func (this *Smime) GetSmime() (res []byte){
	//参考 openssl 实现
	if this.Pkcs7 == nil || this.Msg == nil {
		return nil
	}
	var resBuff bytes.Buffer

	mime_eol := "\n"
	if false {mime_eol = "\r\n"}

	mime_prefix := "application/x-pkcs7-"
	if false {mime_prefix = "application/pkcs7-"}

	bound := make([]byte,32)
	_, err := rand.Read(bound)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for i := 0; i < 32; i++ {
		b := bound[i] & 0xf;
		if (b < 10){
			b += '0'; //48
		}else{
			b += 'A' - 10;
		}
		bound[i] = b
	}
	this.Bound = string(bound)

	resBuff.Write([]byte(`MIME-Version: 1.0` + mime_eol))
	resBuff.Write([]byte(`Content-Type: multipart/signed; protocol="`+ mime_prefix +`signature"; micalg="` + "sha256"))
	resBuff.Write([]byte(`"; boundary="----`+ this.Bound +`"` + mime_eol + mime_eol))
	resBuff.Write([]byte(`This is an S/MIME signed message` + mime_eol + mime_eol))
	resBuff.Write([]byte(`------` + this.Bound + mime_eol))
	resBuff.Write(this.Msg)		//写入消息主体
	resBuff.Write([]byte(mime_eol + `------` + this.Bound + mime_eol))

	//Header for signature
	resBuff.Write([]byte(`Content-Type: `+mime_prefix+`signature; name="smime.p7s"` + mime_eol))
	resBuff.Write([]byte(`Content-Transfer-Encoding: base64` + mime_eol))
	resBuff.Write([]byte(`Content-Disposition: attachment; filename="smime.p7s"` + mime_eol + mime_eol))
	//输出pkcs7主体
	breaker := LineBreaker{out:&resBuff}
	breaker.Handle([]byte(base64.StdEncoding.EncodeToString(this.Pkcs7)),64)
	resBuff.Write([]byte(mime_eol + `------`+ this.Bound +`--` + mime_eol + mime_eol))
	res = resBuff.Bytes()
	return
}

