import { SiteHeader } from "@/components/SiteHeader";
import { SiteFooter } from "@/components/SiteFooter";
import { PromoRibbon } from "@/components/PromoRibbon";
import { CartProvider } from "@/lib/cart";

export default function SiteLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <CartProvider>
      <SiteHeader />
      <main className="flex-1">{children}</main>
      <SiteFooter />
      {/* Spacer so the fixed bottom ribbon doesn't cover footer content */}
      <div className="h-16" aria-hidden />
      <PromoRibbon />
    </CartProvider>
  );
}
