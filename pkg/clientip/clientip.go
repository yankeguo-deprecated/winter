package clientip

import (
	"github.com/guoyk93/winter"
	"net"
	"strings"
)

// Get extract client ip from 'X-Forwarded-For' header
func Get(c winter.Context) (o string) {
	xff := strings.Join(c.Req().Header.Values("X-Forwarded-For"), ",")

	if xff != "" {
		for _, item := range strings.Split(xff, ",") {
			item = strings.TrimSpace(item)
			if ip := net.ParseIP(item); ip != nil && ip.IsGlobalUnicast() && !ip.IsPrivate() {
				o = item
				break
			}
		}
	}

	if o == "" {
		o, _, _ = net.SplitHostPort(c.Req().RemoteAddr)
	}

	return
}
