import { Route, Routes } from 'react-router-dom'
import { useAuth } from 'react-oidc-context'
import UserHome from './UserHome.tsx'
import MonthlyRecordPage from './pages/MonthlyRecordPage.tsx'
import AssetPage from './pages/AssetPage.tsx'
import AssetCategoryPage from './pages/AssetCategoryPage.tsx'
import TagPage from './pages/TagPage.tsx'
import RecordDetailPage from './pages/RecordDetailPage.tsx'


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
        <Route path="/record/monthly/:year/:month" element={<MonthlyRecordPage />}></Route>
        <Route path="/record/:recordId" element={<RecordDetailPage />}></Route>
        <Route path="/asset" element={<AssetPage />}></Route>
        <Route path="/asset/category" element={<AssetCategoryPage />}></Route>
        <Route path="/tag" element={<TagPage />}></Route>
      </Routes>
      <button onClick={() => signOutRedirect()}>Log Out</button>
    </div>
  )
}

export default App
