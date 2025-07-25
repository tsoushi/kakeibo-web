import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { Client, Provider, cacheExchange, fetchExchange } from 'urql'
import { authExchange } from '@urql/exchange-auth'
import { AuthProvider } from 'react-oidc-context'
import type { AuthProviderProps } from 'react-oidc-context'
import type { AuthConfig, AuthUtilities } from '@urql/exchange-auth'
import { User } from 'oidc-client-ts'
import { AuthGuard } from './AuthGuard.tsx'

const cognitoAuthConfig: AuthProviderProps = {
  authority: import.meta.env.VITE_COGNITO_AUTHORITY,
  client_id: import.meta.env.VITE_COGNITO_CLIENT_ID,
  redirect_uri: import.meta.env.VITE_FRONT_URL,
  response_type: "code",
  scope: "email openid phone",
}

const getCognitoUser = () => {
  const oidcStorage = sessionStorage.getItem(`oidc.user:${cognitoAuthConfig.authority}:${cognitoAuthConfig.client_id}`)

  if (!oidcStorage) {
    return null
  }

  return User.fromStorageString(oidcStorage)
}

const authConfig = async (utils: AuthUtilities): Promise<AuthConfig> => {
  return {
    addAuthToOperation(operation) {
      const user = getCognitoUser()
      if (!user) {
        console.log('not authenticated')
        return operation
      }

      const op = utils.appendHeaders(operation, {
        Authorization: `Bearer ${user.access_token}`,
      })
      return op
    },
    didAuthError(error) {
      return error.graphQLErrors.some(e => e.extensions?.code === 'UNAUTHENTICATED')
    },
    async refreshAuth() {
      console.log('Refreshing auth...')
    },
    willAuthError(operation) {
      const user = getCognitoUser()

      if (!user) {
        return true
      }
      return false
    }
  }
}

const client = new Client({
  url: import.meta.env.VITE_GRAPHQL_SERVER_URL,
  exchanges: [
    cacheExchange,
    authExchange(authConfig),
    fetchExchange,
  ],
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider {...cognitoAuthConfig}>
      <Provider value={client}>
        <AuthGuard>
          <App />
        </AuthGuard>
      </Provider>
    </AuthProvider>
  </StrictMode>,
)
