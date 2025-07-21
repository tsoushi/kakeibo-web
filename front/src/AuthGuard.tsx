import type { ReactNode, FC } from "react";
import { useAuth } from "react-oidc-context";


export const AuthGuard: FC<{children?: ReactNode}> = ({ children }) => {
  const auth = useAuth()

  const signOutRedirect = () => {
    const clientId = import.meta.env.VITE_COGNITO_CLIENT_ID;
    const logoutUri = "http://localhost:5173/";
    const cognitoDomain = import.meta.env.VITE_COGNITO_DOMAIN;
    window.location.href = `${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(logoutUri)}`;
  };

  if (!auth.isAuthenticated) {
    return <div>
      <p>Please log in to view user information.</p>
      <button onClick={() => auth.signinRedirect()}>Log In</button>
      <button onClick={() => signOutRedirect()}>Log Out</button>
    </div>
  }
  if (auth.isLoading) {
    return <div>Loading authentication...</div>
  }

  console.log('User is authenticated:', auth.user)

  return (
    <div>
      {children}
      <button onClick={() => auth.signoutRedirect()}>Log Out</button>
    </div>
  )
}