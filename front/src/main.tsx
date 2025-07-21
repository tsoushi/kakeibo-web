import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { Client, Provider, cacheExchange, fetchExchange } from 'urql'
import { authExchange } from '@urql/exchange-auth'
import type { AuthConfig, AuthUtilities } from '@urql/exchange-auth'
import { AuthProvider, type AuthProviderProps } from 'react-oidc-context'
import { AuthGate } from './AuthGate.tsx'

const authConfig = async (utils: AuthUtilities): Promise<AuthConfig> => {
  const value = localStorage.getItem("Amazon.AWS.Cognito.ContextData.LS_UBID")

  return {
    addAuthToOperation(operation) {
      const op = utils.appendHeaders(operation, {
        "Debug-User-Name": import.meta.env.VITE_DEBUG_USER_NAME,
        "Debug-User-Password": import.meta.env.VITE_DEBUG_USER_PASSWORD,
        "access-token": value || '',
      })
      return op
    },
    didAuthError(error) {
      return error.graphQLErrors.some(e => e.extensions?.code === 'UNAUTHENTICATED')
    },
    async refreshAuth() {
    },
    willAuthError(operation) {
      return false
    }
  }
}

const cognitoAuthConfig = {
  authority: "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_4WRSPQRwE",
  client_id: "319he56kb8i60id536gcr0sirh",
  redirect_uri: "http://localhost:5173/",
  response_type: "code",
  scope: "email openid phone",
};

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
    <Provider value={client}>
      <AuthProvider {...cognitoAuthConfig}>
        <AuthGate>
          <App /> 
        </AuthGate>
      </AuthProvider>
    </Provider>
  </StrictMode>,
)
