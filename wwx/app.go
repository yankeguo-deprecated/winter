package wwx

import (
	"github.com/guoyk93/winter"
	"log"
	"net/http"
)

type app struct {
	opt *options
}

func (a *app) handleValidation(c winter.Context) {
	req := winter.Bind[struct {
		Signature string `json:"signature"`
		Timestamp int64  `json:"timestamp,string"`
		Nonce     string `json:"nonce"`
		Echostr   string `json:"echostr"`
	}](c)

	c.Text(req.Echostr)
}

func (a *app) handleIncoming(c winter.Context) {
	req := winter.Bind[map[string]any](c)

	log.Println(req)

	c.Text("success")
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
