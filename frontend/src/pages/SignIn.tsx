import { FormEvent, useState } from 'react';
import { useSession } from '../auth/SessionProvider';

export function SignIn() {
  const { login, loading, error } = useSession();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const canSubmit = email.trim() !== '' && password !== '' && !loading;

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (canSubmit) {
      await login(email, password);
    }
  }

  return (
    <main className="signInPage">
      <form className="signInPanel" onSubmit={submit}>
        <div>
          <h1>CRM System</h1>
          <p>Sign in to continue</p>
        </div>
        <label>
          Email
          <input autoFocus type="email" value={email} onChange={(event) => setEmail(event.target.value)} />
        </label>
        <label>
          Password
          <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
        </label>
        {error && <div role="alert" className="error">{error}</div>}
        <button className="primaryButton" type="submit" disabled={!canSubmit}>
          {loading ? 'Signing in' : 'Sign in'}
        </button>
      </form>
    </main>
  );
}
