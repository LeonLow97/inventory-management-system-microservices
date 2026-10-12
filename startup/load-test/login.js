import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = __ENV.K6_BASE_URL || 'http://localhost:8080';
const USERS = JSON.parse(__ENV.K6_USERS || '[]');
const REQUESTS_PER_USER = 200;

if (!Array.isArray(USERS) || USERS.length !== 5) {
  throw new Error('K6_USERS must be a JSON array containing five user objects.');
}

const usersAreValid = USERS.every((user) =>
  user &&
  typeof user.email === 'string' &&
  user.email.trim().length > 0 &&
  typeof user.password === 'string' &&
  user.password.length > 0,
);

if (!usersAreValid) {
  throw new Error('K6_USERS must contain five user objects with an email and password each.');
}

const emails = USERS.map((user) => user.email.trim().toLowerCase());

if (new Set(emails).size !== USERS.length) {
  throw new Error('K6_USERS must contain five unique email addresses.');
}

const params = {
  headers: { 'Content-Type': 'application/json' },
};

export const options = {
  scenarios: {
    concurrent_login_burst: {
      executor: 'per-vu-iterations',
      vus: USERS.length * REQUESTS_PER_USER,
      iterations: 1,
      maxDuration: '2m',
    },
  },
};

export default function () {
  const userIndex = Math.floor((__VU - 1) / REQUESTS_PER_USER);
  const user = USERS[userIndex];
  const payload = JSON.stringify({ email: user.email, password: user.password });
  const res = http.post(`${BASE_URL}/api/v1/login`, payload, params);

  check(res, {
    'status is 200 or rate-limited (429)': (r) => r.status === 200 || r.status === 429,
    'response time < 1s': (r) => r.timings.duration < 1000,
  });
}
