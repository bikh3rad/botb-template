import type { SVGProps } from "react";

/** BOTB orange wordmark logo. Source SVG lives at /images/ui/logo.svg. */
export function Logo({ className }: { className?: string }) {
  // eslint-disable-next-line @next/next/no-img-element
  return <img src="/images/ui/logo.svg" alt="BOTB" className={className} />;
}

export function CartIcon({ className, ...props }: SVGProps<SVGSVGElement> & { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="26"
      viewBox="0 0 24 26"
      fill="none"
      className={className}
      {...props}
    >
      <path
        d="M5.99805 17.4805C6.52849 17.4805 7.03719 17.6912 7.41226 18.0662C7.78734 18.4414 7.99805 18.95 7.99805 19.4804C7.99805 20.011 7.78734 20.5196 7.41226 20.8946C7.03719 21.2698 6.52849 21.4805 5.99805 21.4805C5.46761 21.4805 4.95891 21.2698 4.58384 20.8946C4.20877 20.5196 3.99805 20.011 3.99805 19.4804C3.99805 18.95 4.20877 18.4414 4.58384 18.0662C4.95891 17.6912 5.46761 17.4805 5.99805 17.4805ZM5.99805 17.4805H16.998M5.99805 17.4805V3.48047H3.99805M16.998 17.4805C17.5284 17.4805 18.0372 17.6912 18.4122 18.0662C18.7874 18.4414 18.9981 18.95 18.9981 19.4804C18.9981 20.011 18.7874 20.5196 18.4122 20.8946C18.0372 21.2698 17.5284 21.4805 16.998 21.4805C16.4676 21.4805 15.959 21.2698 15.5838 20.8946C15.2087 20.5196 14.9981 20.011 14.9981 19.4804C14.9981 18.95 15.2087 18.4414 15.5838 18.0662C15.959 17.6912 16.4676 17.4805 16.998 17.4805ZM5.99805 5.48047L19.998 6.48047L18.9981 13.4804H5.99805"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

export function ChevronLeftIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" className={className} {...props}>
      <path d="M15 18l-6-6 6-6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

export function ChevronRightIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" className={className} {...props}>
      <path d="M9 18l6-6-6-6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

export function ChevronDownIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" className={className} {...props}>
      <path d="M6 9l6 6 6-6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

export function MenuIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" className={className} {...props}>
      <path d="M4 6h16M4 12h16M4 18h16" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
  );
}

export function UserIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" className={className} {...props}>
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="1.5" />
      <circle cx="12" cy="10" r="3" stroke="currentColor" strokeWidth="1.5" />
      <path d="M6.5 19a6 6 0 0 1 11 0" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
    </svg>
  );
}

export function CloseIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" className={className} {...props}>
      <path d="M6 6l12 12M18 6L6 18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
  );
}

export function TicketIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className={className} {...props}>
      <path d="M3 8a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v2a2 2 0 0 0 0 4v2a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-2a2 2 0 0 0 0-4V8z" />
    </svg>
  );
}

export function LightningIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className={className} {...props}>
      <path d="M13 2L4.5 13.5H11l-1 8.5L19.5 10H13l0-8z" />
    </svg>
  );
}
