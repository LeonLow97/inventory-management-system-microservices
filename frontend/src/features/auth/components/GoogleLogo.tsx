import SvgIcon, { type SvgIconProps } from '@mui/material/SvgIcon';

export function GoogleLogo(props: SvgIconProps) {
  return (
    <SvgIcon viewBox="0 0 24 24" aria-hidden="true" {...props}>
      <path
        fill="#4285F4"
        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.27-4.74 3.27-8.1z"
      />
      <path
        fill="#34A853"
        d="M12 23c2.97 0 5.46-.98 7.28-2.65l-3.57-2.77c-.99.66-2.25 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84A11 11 0 0 0 12 23z"
      />
      <path
        fill="#FBBC05"
        d="M5.84 14.11c-.22-.66-.35-1.37-.35-2.11s.13-1.45.35-2.11V7.05H2.18A11 11 0 0 0 1 12c0 1.78.43 3.47 1.18 4.95l3.66-2.84z"
      />
      <path
        fill="#EA4335"
        d="M12 5.25c1.62 0 3.06.56 4.21 1.64l3.15-3.15A10.5 10.5 0 0 0 12 1a11 11 0 0 0-9.82 6.05l3.66 2.84C6.71 7.18 9.14 5.25 12 5.25z"
      />
    </SvgIcon>
  );
}
