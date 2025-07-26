package cognito

import "fmt"

type Config struct {
	Region     string
	UserPoolID string
}

func (c Config) Issuer() string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", c.Region, c.UserPoolID)
}

func (c Config) JWKSURL() string {
	return c.Issuer() + "/.well-known/jwks.json"
}
