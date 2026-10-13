import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import '@fontsource/inter/latin-400.css';
import '@fontsource/inter/latin-500.css';
import '@fontsource/inter/latin-600.css';
import '@fontsource/inter/latin-700.css';
import '@fontsource/manrope/latin-600.css';
import '@fontsource/manrope/latin-700.css';
import '@fontsource/manrope/latin-800.css';
import { App } from './app/App';
import { AppProviders } from './app/AppProviders';
import { ErrorBoundary } from './app/ErrorBoundary';

const root = document.getElementById('root');
if (!root) throw new Error('The application root element is missing.');

createRoot(root).render(
  <StrictMode>
    <AppProviders>
      <ErrorBoundary>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </ErrorBoundary>
    </AppProviders>
  </StrictMode>,
);
