package wjwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter/wjwk"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	rk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	k, err := jwk.FromRaw(rk)
	require.NoError(t, err)
	k.Set(jwk.KeyUsageKey, jwk.ForSignature)
	k.Set(jwk.AlgorithmKey, jwa.KeyAlgorithmFrom(jwa.RS256))
	ctx := context.Background()
	ctx = context.WithValue(ctx, wjwk.Ext.Instance("test").ContextKey(), k)
	ctx = context.WithValue(ctx, Ext.Instance("test").ContextKey(), Ext.Options(
		WithIssuer("test-issuer"),
		WithJWKKey("test"),
	))
	var s string
	func() {
		defer rg.Guard(&err)
		s = Sign(ctx, func(b *jwt.Builder) *jwt.Builder {
			return b.Subject("test-sub")
		}, "test")
	}()
	require.NoError(t, err)
	require.Equal(t, 3, len(strings.Split(s, ".")))
}
