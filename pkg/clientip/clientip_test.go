package clientip

import (
	"github.com/guoyk93/winter"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {

	a := winter.New()
	a.HandleFunc("/test", func(c winter.Context) {
		c.Text(Get(c))
	})

	s := httptest.NewServer(a)
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL+"/test", nil)
	require.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "10.10.10.10,203.0.113.195,2001:db8:85a3:8d3:1319:8a2e:370:7348,150.172.238.178")
	res, err := s.Client().Do(req)
	require.NoError(t, err)
	buf, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "203.0.113.195", string(buf))
}
