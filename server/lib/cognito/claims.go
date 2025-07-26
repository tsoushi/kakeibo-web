package cognito

import "github.com/golang-jwt/jwt/v5"

type CognitoClaims struct {
	Username string `json:"username"`
	TokenUse string `json:"token_use"` // "access" or "id"
	Scope    string `json:"scope,omitempty"`
	jwt.RegisteredClaims
}
