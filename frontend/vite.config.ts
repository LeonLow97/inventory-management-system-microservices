import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';
import { loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const apiProxy = {
    '/api': {
      target: env.API_PROXY_TARGET || 'http://127.0.0.1:8080',
      changeOrigin: true,
      // A host-only cookie is bound to the frontend origin by the browser.
      // This handles the different gateway cookie domains in local/Docker/Kind.
      cookieDomainRewrite: '',
    },
  };

  return {
    plugins: [react()],
    server: { port: 3000, strictPort: true, proxy: apiProxy },
    preview: { proxy: apiProxy },
    test: {
      environment: 'jsdom',
      setupFiles: ['./src/test/setup.ts'],
      environmentOptions: { jsdom: { url: 'http://localhost:3000' } },
      clearMocks: true,
      restoreMocks: true,
    },
  };
});
