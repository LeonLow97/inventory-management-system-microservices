import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { Brand } from './Brand';

type PageStateProps = {
  title: string;
  description?: string;
  loading?: boolean;
  onRetry?: () => void;
};

export function PageState({ title, description, loading = false, onRetry }: PageStateProps) {
  return (
    <Box component="main" sx={{ minHeight: '100dvh', display: 'grid', placeItems: 'center', p: 3 }}>
      <Stack
        spacing={3}
        sx={{ alignItems: 'center', textAlign: 'center' }}
        role={loading ? 'status' : 'alert'}
        aria-live="polite"
      >
        <Brand />
        {loading && <CircularProgress size={30} aria-label="Loading" />}
        <Box>
          <Typography component="h1" variant="h3" sx={{ fontSize: 30 }}>
            {title}
          </Typography>
          {description && (
            <Typography color="text.secondary" sx={{ mt: 1, maxWidth: 360 }}>
              {description}
            </Typography>
          )}
        </Box>
        {onRetry && (
          <Button variant="contained" onClick={onRetry}>
            Try again
          </Button>
        )}
      </Stack>
    </Box>
  );
}
