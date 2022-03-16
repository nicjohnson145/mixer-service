package authtest

import (
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	DefaultUsername = "foobar"
)

type AuthOpts struct {
	Username *string
}

func AuthenticatedRequest(t *testing.T, r *http.Request, opts AuthOpts) {
	if opts.Username == nil {
		opts.Username = to.StringPtr(DefaultUsername)
	}

	token, err := auth.GenerateAccessToken(auth.TokenInputs{
		Username: *opts.Username,
	})
	require.NoError(t, err)

	r.Header.Set(auth.AuthenticationHeader, token)
}
