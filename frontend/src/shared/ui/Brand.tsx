import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

export function Brand({ inverse = false }: { inverse?: boolean }) {
  return (
    <Stack direction="row" spacing={1.25} sx={{ alignItems: 'center' }}>
      <Box
        aria-hidden="true"
        sx={{
          width: 38,
          height: 38,
          display: 'grid',
          placeItems: 'center',
          borderRadius: '11px',
          background: inverse ? 'rgba(255,255,255,0.12)' : 'primary.main',
          color: 'white',
          border: inverse ? '1px solid rgba(255,255,255,0.16)' : undefined,
        }}
      >
        <LockOutlinedIcon sx={{ fontSize: 22 }} />
      </Box>
      <Typography
        sx={{
          fontFamily: 'Manrope, sans-serif',
          fontWeight: 800,
          fontSize: 25,
          letterSpacing: '-0.045em',
          color: inverse ? 'white' : 'text.primary',
        }}
      >
        authify
        <Box component="span" sx={{ color: inverse ? '#85aaff' : 'primary.main' }}>
          .
        </Box>
      </Typography>
    </Stack>
  );
}
