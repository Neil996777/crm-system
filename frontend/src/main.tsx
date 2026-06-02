import React from 'react';
import { createRoot } from 'react-dom/client';
import { SessionProvider } from './auth/SessionProvider';
import { Shell } from './app/Shell';
import './styles.css';

createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <SessionProvider>
      <Shell />
    </SessionProvider>
  </React.StrictMode>
);
