package wwx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/pkg/pkcs7pad"
	"net/http"
	"sort"
	"strings"
)

type EncryptedRequest struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	OpenID       string `json:"openid"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Body         []byte `json:"body"`
}

type EncryptedRequestData struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string   `xml:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"`
}

type EncryptedResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName CDATA    `xml:"ToUserName"`
	Encrypt    CDATA    `xml:"Encrypt"`
}

type DecryptAESOptions struct {
	Token          string
	AESKey         string
	AppID          string
	SkipValidation bool
}

type EncryptAESOptions struct {
	Token      string
	AESKey     string
	AppID      string
	ToUserName string
}

func EncryptAES(buf []byte, opts EncryptAESOptions) (res EncryptedResponse, err error) {
	defer rg.Guard(&err)

	// rebuild buf with prefix and suffix
	{
		bufPfx := make([]byte, 20, 20)

		// build prefix with rand and len
		{
			bufRnd := make([]byte, 8, 8)
			rg.Must(rand.Read(bufRnd))
			hex.Encode(bufPfx, bufRnd)

			binary.BigEndian.PutUint32(bufPfx[16:], uint32(len(buf)))
		}

		buf = append(append(bufPfx, buf...), []byte(opts.AppID)...)
	}

	// pkcs7
	{
		buf = pkcs7pad.Encode(buf, 32)
	}

	// encrypt
	{
		key := rg.Must(base64.RawURLEncoding.DecodeString(opts.AESKey))
		cip := rg.Must(aes.NewCipher(key))
		enc := cipher.NewCBCEncrypter(cip, key[:16])
		enc.CryptBlocks(buf, buf)
	}

	// assign values
	{
		res.ToUserName.Value = opts.ToUserName
		res.Encrypt.Value = base64.StdEncoding.EncodeToString(buf)
	}

	return
}

func DecryptAES(req EncryptedRequest, opts DecryptAESOptions) (buf []byte, err error) {
	defer func() {
		if err == nil {
			return
		}
		err = winter.NewHaltError(err, winter.HaltWithStatusCode(http.StatusBadRequest))
	}()
	defer rg.Guard(&err)

	var data EncryptedRequestData
	rg.Must0(xml.Unmarshal(req.Body, &data))

	// validate signature
	if opts.SkipValidation {
		goto doneValidation
	}

	// validation
	{
		// validate encrypt_type
		if req.EncryptType != encryptTypeAES {
			err = errors.New("wwx.DecryptAES: invalid 'ecrypt_type', must be 'aes'")
			return
		}

		// build to sign
		components := []string{
			opts.Token,
			req.Timestamp,
			req.Nonce,
			data.Encrypt,
		}
		sort.Strings(components)

		// digest
		h := sha1.New()
		h.Write([]byte(strings.Join(components, "")))
		d := hex.EncodeToString(h.Sum(nil))

		// compare
		if strings.ToLower(req.MsgSignature) != strings.ToLower(d) {
			err = errors.New("wwx.DecryptAES: invalid 'msg_signature'")
			return
		}
	}

doneValidation:

	// decode base64
	{
		buf = rg.Must(base64.StdEncoding.DecodeString(data.Encrypt))
	}

	// decrypt
	{

		key := rg.Must(base64.RawURLEncoding.DecodeString(opts.AESKey))
		cip := rg.Must(aes.NewCipher(key))
		dec := cipher.NewCBCDecrypter(cip, key[:16])
		dec.CryptBlocks(buf, buf)
	}

	// pkcs#7
	{
		buf = rg.Must(pkcs7pad.Decode(buf))
	}

	// check and remove suffix
	{
		sfx := []byte(opts.AppID)
		if !bytes.HasSuffix(buf, sfx) {
			err = errors.New("wwx.DecryptAES: missing APP_ID suffix in decrypted data")
			return
		}
		buf = bytes.TrimSuffix(buf, sfx)
	}

	// check and remove prefix
	{
		if len(buf) < 20 {
			err = errors.New("wwx.DecryptAES: decrypted data is too short")
			return
		}

		bufLen := int(binary.BigEndian.Uint32(buf[16:20]))

		buf = buf[20:]

		if len(buf) != bufLen {
			err = errors.New("wwx.DecryptAES: invalid length prefix in decrypted data")
		}
	}

	return
}
