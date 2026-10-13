import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getCurrentUser, login, logout } from '../api/auth-api';

const currentUserKey = ['auth', 'current-user'] as const;

export function useCurrentUser() {
  return useQuery({
    queryKey: currentUserKey,
    queryFn: ({ signal }) => getCurrentUser(signal),
    retry: false,
    staleTime: 30_000,
    refetchInterval: (query) => (query.state.data ? 60_000 : false),
    refetchOnWindowFocus: 'always',
  });
}

export function useLogin() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: login,
    onSuccess: async (user) => {
      // Prevent an earlier unauthenticated profile request from replacing this login.
      await client.cancelQueries({ queryKey: currentUserKey });
      client.setQueryData(currentUserKey, user);
    },
  });
}

export function useLogout() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: logout,
    onSuccess: async () => {
      await client.cancelQueries({ queryKey: currentUserKey });
      client.setQueryData(currentUserKey, null);
    },
  });
}
