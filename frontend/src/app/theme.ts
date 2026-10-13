import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  palette: {
    primary: { main: '#2459e0', dark: '#1743b7', light: '#eaf0ff' },
    background: { default: '#f7f9fe', paper: '#ffffff' },
    text: { primary: '#111827', secondary: '#627089' },
    divider: '#e6ebf4',
  },
  typography: {
    fontFamily: 'Inter, sans-serif',
    body1: { fontSize: '1.125rem' },
    body2: { fontSize: '1rem' },
    h1: { fontFamily: 'Manrope, sans-serif', fontWeight: 800, letterSpacing: '-0.055em' },
    h2: { fontFamily: 'Manrope, sans-serif', fontWeight: 800, letterSpacing: '-0.045em' },
    h3: { fontFamily: 'Manrope, sans-serif', fontWeight: 700, letterSpacing: '-0.035em' },
    button: { fontSize: '1rem', fontWeight: 600, textTransform: 'none', letterSpacing: 0 },
    overline: { fontSize: '0.75rem', fontWeight: 700, letterSpacing: '0.16em' },
  },
  shape: { borderRadius: 12 },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        'html, body, #root': { minHeight: '100%', margin: 0 },
        body: { minWidth: 320 },
        '*': { boxSizing: 'border-box' },
        '::selection': { background: '#dbe6ff' },
        '@media (prefers-reduced-motion: reduce)': {
          '*, *::before, *::after': {
            animationDuration: '0.01ms !important',
            animationIterationCount: '1 !important',
            transitionDuration: '0.01ms !important',
            scrollBehavior: 'auto !important',
          },
        },
      },
    },
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: { borderRadius: 10, padding: '12px 20px' },
        sizeLarge: { fontSize: '1.0625rem' },
        sizeSmall: { fontSize: '0.875rem' },
      },
    },
    MuiInputBase: { styleOverrides: { root: { fontSize: '1.0625rem' } } },
    MuiInputLabel: { styleOverrides: { root: { fontSize: '1.0625rem' } } },
    MuiFormHelperText: { styleOverrides: { root: { fontSize: '0.875rem' } } },
    MuiOutlinedInput: {
      styleOverrides: {
        root: { borderRadius: 10, background: '#ffffff' },
        notchedOutline: { borderColor: '#d9e0ec' },
      },
    },
    MuiAlert: { styleOverrides: { root: { borderRadius: 10 } } },
  },
});
