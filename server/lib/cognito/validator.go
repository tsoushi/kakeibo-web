package cognito

import (
	"context"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/xerrors"
)

type Validator struct {
	cfg     Config
	keyfunc keyfunc.Keyfunc
}

func NewCognitoValidator(ctx context.Context, cfg Config) (*Validator, error) {
	kf, err := keyfunc.NewDefaultCtx(ctx, []string{cfg.JWKSURL()})
	if err != nil {
		return nil, xerrors.Errorf("failed to create keyfunc: %w", err)
	}

	return &Validator{cfg: cfg, keyfunc: kf}, nil
}

func (v *Validator) ValidateToken(tokenString string) (*CognitoClaims, error) {
	var claims CognitoClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, v.keyfunc.Keyfunc,
		jwt.WithIssuer(v.cfg.Issuer()),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to validate token: %w", err)
	}

	return &claims, nil
}
