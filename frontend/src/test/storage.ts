import { vi } from 'vitest';

function storageMock() {
  return {
    length: 0,
    clear: vi.fn(),
    getItem: vi.fn(() => null),
    key: vi.fn(() => null),
    removeItem: vi.fn(),
    setItem: vi.fn<(key: string, value: string) => void>(),
  } satisfies Storage;
}

export const storage = { local: storageMock(), session: storageMock() };
