package wwx

import (
	"encoding/xml"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"net/http"
	"time"
)

const (
	encryptTypeAES = "aes"
)

type app struct {
	opt *options
}

type ValidationRequest struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp,string"`
	Nonce     string `json:"nonce"`
	Echostr   string `json:"echostr"`
}

func (a *app) handleValidation(c winter.Context) {
	req := winter.Bind[ValidationRequest](c)
	c.Text(req.Echostr)
}

func (a *app) handleIncoming(c winter.Context) {
	encReq := winter.Bind[EncryptedRequest](c)

	buf := rg.Must(DecryptAES(encReq, DecryptAESOptions{
		Token:          a.opt.token,
		AESKey:         a.opt.aesKey,
		AppID:          a.opt.appID,
		SkipValidation: a.opt.skipValidation,
	}))

	var req Request
	rg.Must0(xml.Unmarshal(buf, &req))

	var h HandlerFunc

	if req.MsgType == TypeEvent && a.opt.evtHandlers != nil {
		h = a.opt.evtHandlers[req.Event]
		if h == nil {
			h = a.opt.evtHandlers[""]
		}
	} else if a.opt.msgHandlers != nil {
		h = a.opt.msgHandlers[req.MsgType]
		if h == nil {
			h = a.opt.msgHandlers[""]
		}
	}

	if h == nil {
		c.Text("success")
		return
	}

	res := &Response{}
	res.ToUserName.Value = req.FromUserName
	res.FromUserName.Value = req.ToUserName
	res.CreateTime = time.Now().Unix()

	wc := &wxContext{c: c, req: req, res: res}

	h(wc)

	if res.Empty || res.MsgType.Value == "" {
		c.Text("success")
		return
	}

	encRes := rg.Must(EncryptAES(
		rg.Must(xml.Marshal(res)),
		EncryptAESOptions{
			Token:      a.opt.token,
			AESKey:     a.opt.aesKey,
			AppID:      a.opt.appID,
			ToUserName: req.FromUserName,
		},
	))

	c.Body("application/xml", rg.Must(xml.Marshal(encRes)))
}

func (a *app) HandleCallback(c winter.Context) {
	if c.Req().Method == "GET" {
		a.handleValidation(c)
	} else if c.Req().Method == "POST" {
		a.handleIncoming(c)
	} else {
		winter.HaltString("unexpected method", winter.HaltWithStatusCode(http.StatusMethodNotAllowed))
	}
	return
}
