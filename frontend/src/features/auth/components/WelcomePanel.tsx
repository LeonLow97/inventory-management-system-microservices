import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import CheckIcon from '@mui/icons-material/Check';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { Brand } from '../../../shared/ui/Brand';

// The artwork is decorative CSS and SVG icons, so the screen reader skips it.
function WelcomeArtwork() {
  return (
    <Box
      aria-hidden="true"
      sx={{
        position: 'relative',
        height: 230,
        width: '100%',
        maxWidth: 390,
        mx: 'auto',
        mt: 3,
        mb: 2,
      }}
    >
      {[210, 290, 370].map((size) => (
        <Box
          key={size}
          sx={{
            position: 'absolute',
            width: size,
            height: size,
            borderRadius: '50%',
            border: '1px solid rgba(163,191,255,0.12)',
            left: '50%',
            top: '50%',
            transform: 'translate(-50%, -50%)',
          }}
        />
      ))}
      <Box
        sx={{
          position: 'absolute',
          width: 215,
          height: 180,
          top: 18,
          left: 'calc(50% - 112px)',
          border: '1px solid rgba(255,255,255,0.24)',
          borderRadius: '24px',
          background: 'linear-gradient(140deg, rgba(121,163,255,0.25), rgba(32,65,150,0.3))',
          boxShadow: '0 24px 60px rgba(0,0,0,0.18)',
          transform: 'rotate(-9deg)',
          display: 'grid',
          placeItems: 'center',
        }}
      >
        <Box
          sx={{
            width: 76,
            height: 76,
            borderRadius: '22px',
            bgcolor: '#8aafff',
            color: '#102960',
            display: 'grid',
            placeItems: 'center',
            boxShadow: '0 12px 35px rgba(0,0,0,0.12)',
          }}
        >
          <LockOutlinedIcon sx={{ fontSize: 40 }} />
        </Box>
      </Box>
      <Box
        sx={{
          position: 'absolute',
          top: 145,
          left: 'calc(50% - 12px)',
          width: 166,
          p: 2,
          borderRadius: '16px',
          bgcolor: 'white',
          color: '#16213e',
          transform: 'rotate(5deg)',
          boxShadow: '0 18px 42px rgba(0,0,0,0.2)',
          display: 'flex',
          gap: 1.25,
          alignItems: 'center',
        }}
      >
        <Box
          sx={{
            width: 32,
            height: 32,
            borderRadius: '50%',
            bgcolor: '#eaf0ff',
            color: '#2459e0',
            display: 'grid',
            placeItems: 'center',
          }}
        >
          <CheckIcon fontSize="small" />
        </Box>
        <Box>
          <Typography sx={{ fontSize: 13, fontWeight: 700 }}>Your own space</Typography>
          <Typography sx={{ fontSize: 11, color: '#627089', mt: 0.25 }}>
            One sign-in away
          </Typography>
        </Box>
      </Box>
      <Box
        sx={{
          position: 'absolute',
          left: 22,
          bottom: 10,
          width: 8,
          height: 8,
          borderRadius: '50%',
          bgcolor: '#91b2ff',
        }}
      />
      <Box
        sx={{
          position: 'absolute',
          right: 12,
          top: 20,
          width: 5,
          height: 5,
          borderRadius: '50%',
          bgcolor: '#91b2ff',
        }}
      />
    </Box>
  );
}

export function WelcomePanel() {
  return (
    <Box
      sx={{
        display: { xs: 'none', md: 'flex' },
        flexDirection: 'column',
        justifyContent: 'space-between',
        position: 'relative',
        overflow: 'hidden',
        p: { md: 5, lg: 6 },
        background:
          'radial-gradient(ellipse at 10% 100%, #173b9a 0%, transparent 65%), linear-gradient(155deg, #0d1b3c 0%, #112b68 100%)',
        color: 'white',
      }}
    >
      <Brand inverse />
      <Box sx={{ py: 3, maxWidth: 460 }}>
        <Typography variant="overline" sx={{ color: '#a3beff' }}>
          YOUR ACCOUNT. YOUR SPACE.
        </Typography>
        <Typography
          component="h2"
          variant="h1"
          sx={{ mt: 2, fontSize: { md: 48, lg: 56 }, lineHeight: 1.08 }}
        >
          Make yourself
          <br />
          <Box component="span" sx={{ color: '#95b6ff' }}>
            at home.
          </Box>
        </Typography>
        <Typography sx={{ mt: 3, maxWidth: 340, color: '#c4d1ed', lineHeight: 1.8, fontSize: 17 }}>
          A familiar place, just one sign-in away.
          <br />
          Welcome to Authify.
        </Typography>
        <WelcomeArtwork />
      </Box>
      <Stack
        direction="row"
        sx={{
          color: '#a7b9de',
          fontSize: 12,
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Box component="span">A simple way in.</Box>
        <Box component="span" sx={{ width: 32, height: 2, bgcolor: '#749eff' }} />
      </Stack>
    </Box>
  );
}
