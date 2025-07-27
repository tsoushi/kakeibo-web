import { Client, cacheExchange, fetchExchange } from 'urql'
import { authExchange } from '@urql/exchange-auth'
import type { AuthConfig, AuthUtilities } from '@urql/exchange-auth'
import { getCognitoUser } from './cognito'

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

export const urqlClient = new Client({
  url: import.meta.env.VITE_GRAPHQL_SERVER_URL,
  exchanges: [
    cacheExchange,
    authExchange(authConfig),
    fetchExchange,
  ],
})