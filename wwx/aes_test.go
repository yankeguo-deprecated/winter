package wwx

import (
	"encoding/xml"
	"github.com/guoyk93/rg"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	raw := []byte(`<xml><ToUserName><![CDATA[gh_606a4d298963]]></ToUserName>
<FromUserName><![CDATA[oGnmQ5ixCo0JKJqAunh6i1Yj4yMA]]></FromUserName>
<CreateTime>1677472749</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[xxxxxxx]]></Content>
<MsgId>24015632708048579</MsgId>
</xml>`)

	out, err := EncryptAES(raw,
		EncryptAESOptions{
			Token:     "fd8H8SOsDI6YS2b",
			AESKey:    "fd8H8SOsDI6YS2bYhG6DEFMdmk2mcdEq8kauDUuUajH",
			AppID:     "wxd044d4b039b49604",
			Nonce:     "1713665194",
			Timestamp: "1677472749",
		},
	)
	require.NoError(t, err)

	dout, err := DecryptAES(
		EncryptedRequest{
			Signature:    "",
			Timestamp:    out.Timestamp,
			Nonce:        out.Nonce.Value,
			OpenID:       "oGnmQ5ixCo0JKJqAunh6i1Yj4yMA",
			EncryptType:  "aes",
			MsgSignature: out.MsgSignature.Value,
			Body:         rg.Must(xml.Marshal(out)),
		},
		DecryptAESOptions{
			Token:  "fd8H8SOsDI6YS2b",
			AESKey: "fd8H8SOsDI6YS2bYhG6DEFMdmk2mcdEq8kauDUuUajH",
			AppID:  "wxd044d4b039b49604",
		},
	)

	require.NoError(t, err)
	require.Equal(t, dout, raw)

}

func TestDecryptAES(t *testing.T) {
	out, err := DecryptAES(
		EncryptedRequest{
			Signature:    "9e8dc80a14399c1d8c68e35e79304bdc33daf576",
			Timestamp:    "1677472749",
			Nonce:        "1713665194",
			OpenID:       "oGnmQ5ixCo0JKJqAunh6i1Yj4yMA",
			EncryptType:  "aes",
			MsgSignature: "183acfa46f169492bbb0256683764a548316d1f5",
			Body: []byte(`<xml>

<ToUserName><![CDATA[gh_606a4d298963]]></ToUserName>

<Encrypt><![CDATA[sr1KNfwyDPftbhaWBjMRMVmBwrCIpkNook+wWqO77m0QCoFKecZ9cnlXi/RVMh304Lto9nQwhmVwWF+l53Mp5TXBMj/c9sQ+MH3BttLdy4bNPlXtk+tRebb/sQqKmK3hhIiK6SmMrzmMrxSXR5X/OWNQKVyFPfheKgpUeHIJR5tvBbi92f0M5ct8ZHb7IxBax5fSZhf8qHxtAwIBsdM41L0SRjjojo5j/DSsakxaOLvd7DLELNGvRJkeZ4oK+UP6F0iZ2Lk+sm232lE2SmUZH22CbD7OwyhxYoeBd/5YJpOwFxDb0/gAdYB/ObX9Fq5SUaMtk9A8QDx1nhh/CIOWcPoiHKtJiZr1HFU9acfHfdXJOuADfLEXsJoOq8jOXgPkosi+gHb4XoidIFv7wnP3jduTfMX5/pOD8hkE4qtybek=]]></Encrypt>

</xml>`),
		},
		DecryptAESOptions{
			Token:  "fd8H8SOsDI6YS2b",
			AESKey: "fd8H8SOsDI6YS2bYhG6DEFMdmk2mcdEq8kauDUuUajH",
			AppID:  "wxd044d4b039b49604",
		},
	)

	require.NoError(t, err)

	require.Equal(t, `<xml><ToUserName><![CDATA[gh_606a4d298963]]></ToUserName>
<FromUserName><![CDATA[oGnmQ5ixCo0JKJqAunh6i1Yj4yMA]]></FromUserName>
<CreateTime>1677472749</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[xxxxxxx]]></Content>
<MsgId>24015632708048579</MsgId>
</xml>`, string(out))
}
