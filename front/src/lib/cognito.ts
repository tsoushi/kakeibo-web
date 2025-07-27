import type { AuthProviderProps } from 'react-oidc-context'
import { User } from 'oidc-client-ts'


export const cognitoAuthConfig: AuthProviderProps = {
  authority: import.meta.env.VITE_COGNITO_AUTHORITY,
  client_id: import.meta.env.VITE_COGNITO_CLIENT_ID,
  redirect_uri: import.meta.env.VITE_FRONT_URL,
  response_type: "code",
  scope: "email openid phone",
}

export const getCognitoUser = () => {
  const oidcStorage = sessionStorage.getItem(`oidc.user:${cognitoAuthConfig.authority}:${cognitoAuthConfig.client_id}`)

  if (!oidcStorage) {
    return null
  }

  return User.fromStorageString(oidcStorage)
}
