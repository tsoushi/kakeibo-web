import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { Client, Provider, cacheExchange, fetchExchange } from 'urql'
import { authExchange } from '@urql/exchange-auth'
import type { AuthConfig, AuthUtilities } from '@urql/exchange-auth'

const authConfig = async (utils: AuthUtilities): Promise<AuthConfig> => {
  // let token = localStorage.getItem('token')

  return {
    addAuthToOperation(operation) {
      const op = utils.appendHeaders(operation, {
        "Debug-User-Name": import.meta.env.VITE_DEBUG_USER_NAME,
        "Debug-User-Password": import.meta.env.VITE_DEBUG_USER_PASSWORD,
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
      <App />
    </Provider>
  </StrictMode>,
)
