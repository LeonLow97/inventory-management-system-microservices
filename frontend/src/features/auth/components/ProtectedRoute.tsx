import { Navigate, Outlet } from 'react-router-dom';
import { PageState } from '../../../shared/ui/PageState';
import { useCurrentUser } from '../hooks/use-auth';

export function ProtectedRoute() {
  const auth = useCurrentUser();
  if (auth.isPending) return <PageState title="Checking your session" loading />;
  if (auth.isError) {
    return (
      <PageState
        title="We couldn’t load your account"
        description="Check your connection and try again."
        onRetry={() => {
          void auth.refetch();
        }}
      />
    );
  }
  if (!auth.data) return <Navigate to="/login" replace />;
  return <Outlet />;
}
