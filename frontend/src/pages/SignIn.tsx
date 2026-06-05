import { FormEvent, useState } from 'react';
import { useSession } from '../auth/SessionProvider';
import { appName } from '../i18n/labels';

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
          <h1>{appName}</h1>
          <p>登录以继续</p>
        </div>
        <label>
          邮箱
          <input autoFocus type="email" value={email} onChange={(event) => setEmail(event.target.value)} />
        </label>
        <label>
          密码
          <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
        </label>
        {error && <div role="alert" className="error">{error}</div>}
        <button className="primaryButton" type="submit" disabled={!canSubmit}>
          {loading ? '登录中' : '登录'}
        </button>
      </form>
    </main>
  );
}
