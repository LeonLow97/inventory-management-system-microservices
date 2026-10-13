import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { useCurrentUser, useLogout } from '../features/auth/hooks/use-auth';

export function HomePage() {
  const { data: user } = useCurrentUser();
  const logout = useLogout();
  if (!user) return null;

  const fullName = [user.firstName, user.lastName].filter(Boolean).join(' ') || 'there';

  return (
    <Box
      component="main"
      sx={{
        minHeight: '100dvh',
        display: 'grid',
        placeItems: 'center',
        position: 'relative',
        px: 3,
        bgcolor: 'background.default',
      }}
    >
      <Button
        variant="text"
        size="small"
        disabled={logout.isPending}
        onClick={() => logout.mutate()}
        sx={{ position: 'absolute', top: 2, right: 2, color: 'text.secondary', fontSize: 15 }}
      >
        {logout.isPending ? 'Logging out…' : 'Log out'}
      </Button>
      <Box sx={{ textAlign: 'center' }}>
        {logout.isError && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {logout.error.message}
          </Alert>
        )}
        <Typography
          component="h1"
          variant="h1"
          sx={{ fontSize: { xs: 40, sm: 56 }, lineHeight: 1.2, overflowWrap: 'anywhere' }}
        >
          Hello,{' '}
          <Box component="span" sx={{ color: 'primary.main' }}>
            {fullName}
          </Box>
          .
        </Typography>
      </Box>
    </Box>
  );
}
