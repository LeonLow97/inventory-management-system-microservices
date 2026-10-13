import { Navigate, useNavigate } from 'react-router-dom';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { LoginForm } from '../features/auth/components/LoginForm';
import { WelcomePanel } from '../features/auth/components/WelcomePanel';
import { useCurrentUser, useLogin } from '../features/auth/hooks/use-auth';
import { Brand } from '../shared/ui/Brand';
import { PageState } from '../shared/ui/PageState';

export function LoginPage() {
  const auth = useCurrentUser();
  const login = useLogin();
  const navigate = useNavigate();

  if (auth.isPending) return <PageState title="Checking your session" loading />;
  if (auth.data) return <Navigate to="/" replace />;

  return (
    <Box
      component="main"
      sx={{
        minHeight: '100dvh',
        display: 'grid',
        gridTemplateColumns: { xs: '1fr', md: 'minmax(0, 2fr) minmax(0, 3fr)' },
        bgcolor: 'background.paper',
      }}
    >
      <WelcomePanel />
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          minHeight: '100dvh',
          px: { xs: 3, sm: 6, md: 7 },
          py: { xs: 3.5, md: 5 },
        }}
      >
        <Box sx={{ display: { xs: 'block', md: 'none' } }}>
          <Brand />
        </Box>
        <Box sx={{ display: { xs: 'none', md: 'block' }, textAlign: 'right' }}>
          <Typography variant="overline" color="text.secondary">
            WELCOME TO AUTHIFY
          </Typography>
        </Box>
        <Box
          sx={{
            width: '100%',
            maxWidth: 410,
            mx: 'auto',
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            py: 7,
          }}
        >
          <Typography variant="overline" color="primary" sx={{ mb: 1.75 }}>
            LET’S GET YOU SIGNED IN
          </Typography>
          <Typography
            component="h1"
            variant="h2"
            sx={{ fontSize: { xs: 40, sm: 48 }, lineHeight: 1.2 }}
          >
            Good to see you.
          </Typography>
          <Typography
            color="text.secondary"
            sx={{ mt: 1.75, mb: 5, fontSize: { xs: 17, sm: 18 }, lineHeight: 1.7 }}
          >
            Enter your email and password to continue
            <br />
            to your account.
          </Typography>
          <LoginForm
            onSubmit={(credentials) =>
              login.mutate(credentials, {
                onSuccess: () => {
                  void navigate('/', { replace: true });
                },
              })
            }
            onEdit={() => {
              if (login.isError) login.reset();
            }}
            pending={login.isPending}
            error={login.error?.message}
          />
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{ textAlign: 'center', mt: 3.5, fontSize: 14 }}
          >
            Your space is ready when you are.
          </Typography>
        </Box>
        <Box sx={{ textAlign: 'center' }}>
          <Typography sx={{ fontSize: 12, color: 'text.secondary' }}>
            © {new Date().getFullYear()} Authify
          </Typography>
        </Box>
      </Box>
    </Box>
  );
}
