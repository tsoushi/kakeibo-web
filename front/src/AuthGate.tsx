// src/AuthGate.tsx
import { type ReactNode} from 'react';
import { useAuth } from 'react-oidc-context';

export function AuthGate({ children }: { children: ReactNode }) {
  const auth = useAuth()

  // ローディング状態
  if (auth.isLoading) {
    return <div>Loading...</div>;
  }

  // エラー状態
  if (auth.error) {
    return <div>Encountering error... {auth.error.message}</div>;
  }

  // 認証済みの場合は App を表示
  if (auth.isAuthenticated) {
    console.log('User is authenticated:', auth.user);
    return <>{children}</>;
  }

  // 未認証の場合はログインボタンを表示
  return (
    <div>
      <button onClick={() => auth.signinRedirect()}>Sign in</button>
      {/* サンプルにあった signOutRedirect も残したいならここで */}
      <button
        onClick={() => {
          const clientId = '319he56kb8i60id536gcr0sirh';
          const logoutUri = '<logout uri>'; // 必要に応じて差し替え
          const cognitoDomain =
            'https://us-east-14wrspqrwe.auth.us-east-1.amazoncognito.com';
          window.location.href = `${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(
            logoutUri,
          )}`;
        }}
      >
        Sign out
      </button>
    </div>
  );
}
