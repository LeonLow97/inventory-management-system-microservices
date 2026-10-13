import { useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import VisibilityOutlinedIcon from '@mui/icons-material/VisibilityOutlined';
import VisibilityOffOutlinedIcon from '@mui/icons-material/VisibilityOffOutlined';
import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import { GoogleLogo } from './GoogleLogo';
import { loginSchema, type LoginCredentials } from '../model/auth-schema';

type LoginFormProps = {
  onSubmit: (credentials: LoginCredentials) => void;
  onEdit: () => void;
  pending: boolean;
  error?: string;
};

export function LoginForm({ onSubmit, onEdit, pending, error }: LoginFormProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [showGoogleSignInNotice, setShowGoogleSignInNotice] = useState(false);
  const { control, handleSubmit } = useForm<LoginCredentials>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
    mode: 'onTouched',
  });

  return (
    <Box
      component="form"
      noValidate
      onSubmit={(event) => {
        void handleSubmit(onSubmit)(event);
      }}
      onChange={onEdit}
      aria-busy={pending}
    >
      <Stack spacing={2.75}>
        {error && (
          <Alert severity="error" role="alert">
            {error}
          </Alert>
        )}
        <Controller
          name="email"
          control={control}
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              {...field}
              inputRef={ref}
              id="login-email"
              label="Email address"
              placeholder="you@example.com"
              type="email"
              autoComplete="username"
              autoFocus
              fullWidth
              disabled={pending}
              error={Boolean(fieldState.error)}
              helperText={fieldState.error?.message}
              slotProps={{ htmlInput: { autoCapitalize: 'none', spellCheck: false } }}
            />
          )}
        />
        <Controller
          name="password"
          control={control}
          render={({ field: { ref, ...field }, fieldState }) => (
            <TextField
              {...field}
              inputRef={ref}
              id="login-password"
              label="Password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="current-password"
              fullWidth
              disabled={pending}
              error={Boolean(fieldState.error)}
              helperText={fieldState.error?.message}
              slotProps={{
                input: {
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        type="button"
                        aria-label={showPassword ? 'Hide password' : 'Show password'}
                        aria-pressed={showPassword}
                        onClick={() => setShowPassword((visible) => !visible)}
                        onMouseDown={(event) => event.preventDefault()}
                        disabled={pending}
                        edge="end"
                      >
                        {showPassword ? (
                          <VisibilityOffOutlinedIcon fontSize="small" />
                        ) : (
                          <VisibilityOutlinedIcon fontSize="small" />
                        )}
                      </IconButton>
                    </InputAdornment>
                  ),
                },
              }}
            />
          )}
        />
        <Button
          type="submit"
          variant="contained"
          size="large"
          disabled={pending}
          endIcon={pending ? <CircularProgress size={18} color="inherit" /> : <ArrowForwardIcon />}
          sx={{ height: 54, mt: 0.5 }}
        >
          {pending ? 'Signing in…' : 'Log in'}
        </Button>
        <Divider sx={{ color: 'text.secondary', fontSize: 12, letterSpacing: '0.12em' }}>
          OR
        </Divider>
        <Button
          type="button"
          variant="outlined"
          size="large"
          fullWidth
          startIcon={<GoogleLogo />}
          onClick={() => setShowGoogleSignInNotice(true)}
          sx={{ height: 54, borderColor: 'divider', color: 'text.primary' }}
        >
          Sign in with Google
        </Button>
        {showGoogleSignInNotice && (
          <Alert severity="info" role="status" onClose={() => setShowGoogleSignInNotice(false)}>
            Google sign-in isn’t configured yet.
          </Alert>
        )}
      </Stack>
    </Box>
  );
}
