import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { AuthProvider } from 'react-oidc-context'
import { AuthGuard } from './AuthGuard.tsx'
import { Provider } from 'urql'
import { urqlClient } from './lib/urql.ts'
import { cognitoAuthConfig } from './lib/cognito.ts'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider {...cognitoAuthConfig}>
      <Provider value={urqlClient}>
        <AuthGuard>
          <App />
        </AuthGuard>
      </Provider>
    </AuthProvider>
  </StrictMode>,
)
