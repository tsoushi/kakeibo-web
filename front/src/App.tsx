import { Route, Routes } from 'react-router-dom'
import { useAuth } from 'react-oidc-context'
import UserHome from './UserHome.tsx'
import RecordPage from './Record.tsx'


function App() {
  const auth = useAuth()

  if (!auth.isAuthenticated && !auth.isLoading) {
    auth.signinRedirect()
    return <div>Redirecting to login...</div>
  }
  if (auth.isLoading) {
    return <div>Loading authentication...</div>
  }

  const signOutRedirect = () => {
    const clientId = import.meta.env.VITE_COGNITO_CLIENT_ID;
    const logoutUri = import.meta.env.VITE_FRONT_URL;
    const cognitoDomain = import.meta.env.VITE_COGNITO_DOMAIN;
    auth.signoutSilent()
    window.location.href = `${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(logoutUri)}`;
  };

  return (
    <div>
      <Routes>
        <Route path="/" element={<UserHome />}></Route>
        <Route path="/record" element={<RecordPage />}></Route>
      </Routes>
      <button onClick={() => signOutRedirect()}>Log Out</button>
    </div>
  )
}

export default App
