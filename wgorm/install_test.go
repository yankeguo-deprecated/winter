package wgorm

import (
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wboot"
	"github.com/guoyk93/winter/wext"
	"github.com/guoyk93/winter/wext/wexttest"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"testing"
)

func TestInstaller(t *testing.T) {
	dsn := wboot.EnvStr("MYSQL_DSN")
	if dsn == "" {
		return
	}

	wexttest.Run(t, ext, wexttest.Options[options, *injected]{
		AltKeys: []string{"test"},
		Options: []Option{
			WithMySQLDSN(dsn),
		},
		OnCheck: func(code int, body []byte) {
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, "gorm.test: OK", string(body))
		},
		Handle: func(c winter.Context, ins wext.Instance[options, *injected]) {
			var v int64
			err := DB(c, "test").Raw("SELECT 12").Scan(&v).Error
			require.NoError(t, err)
			require.Equal(t, int64(12), v)
			c.Text(strconv.FormatInt(v, 10))
		},
		OnHandle: func(code int, body []byte) {
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, "12", string(body))
		},
	})
}
