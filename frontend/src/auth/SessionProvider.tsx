import { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { CurrentUser, currentUser, signIn, signOut } from '../api/auth';
import { ApiError } from '../api/client';

type SessionState = {
  user: CurrentUser | null;
  loading: boolean;
  error: string;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

const SessionContext = createContext<SessionState | null>(null);

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    currentUser()
      .then((response) => setUser(response.user))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  const value = useMemo<SessionState>(
    () => ({
      user,
      loading,
      error,
      login: async (email: string, password: string) => {
        setLoading(true);
        setError('');
        try {
          const response = await signIn(email, password);
          setUser(response.user);
        } catch (caught) {
          const apiError = caught as ApiError;
          setError(apiError.safeMessage || '认证失败。');
          setUser(null);
        } finally {
          setLoading(false);
        }
      },
      logout: async () => {
        setLoading(true);
        await signOut();
        setUser(null);
        setLoading(false);
      }
    }),
    [error, loading, user]
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession() {
  const value = useContext(SessionContext);
  if (!value) {
    throw new Error('SessionProvider is missing');
  }
  return value;
}
