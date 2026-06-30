import type { Metadata } from "next";
import { Jost, Roboto, Josefin_Sans, Geist_Mono } from "next/font/google";
import "./globals.css";

const jost = Jost({
  variable: "--font-jost",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  display: "swap",
});

const roboto = Roboto({
  variable: "--font-roboto",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  display: "swap",
});

const josefin = Josefin_Sans({
  variable: "--font-josefin",
  subsets: ["latin"],
  weight: ["400", "600", "700"],
  display: "swap",
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL("https://www.botb.com"),
  title: "#1 Online Car Competitions UK | Win A Car Now | Real Winners. Every Day.",
  description:
    "We're the leading online car competitions company in the UK. Tickets from just 50p to win a car, cash, tech, holidays and more.",
  openGraph: {
    title: "BOTB — Real winners. Every day.",
    description:
      "We're the leading online car competitions company in the UK. Tickets from just 50p to win a car, cash, tech, holidays and more.",
    siteName: "BOTB",
    images: ["/seo/og-image.png"],
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${jost.variable} ${roboto.variable} ${josefin.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col bg-white text-botb-text">{children}</body>
    </html>
  );
}
