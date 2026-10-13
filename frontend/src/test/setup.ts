import '@testing-library/jest-dom/vitest';
import { afterAll, afterEach, beforeAll, vi } from 'vitest';
import { cleanup } from '@testing-library/react';
import { server } from './server';
import { storage } from './storage';

// Node 25 exposes its own storage globals. Use browser storage test doubles
// so the tests can detect accidental persistence on every supported Node version.
vi.stubGlobal('localStorage', storage.local);
vi.stubGlobal('sessionStorage', storage.session);

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => {
  cleanup();
  server.resetHandlers();
});
afterAll(() => server.close());
