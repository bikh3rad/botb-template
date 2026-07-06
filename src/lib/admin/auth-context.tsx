"use client";

import * as React from "react";
import { useRouter } from "next/navigation";

export interface AdminProfile {
  id: string;
  name: string;
  email: string;
  role: string;
}

interface AuthState {
  admin: AdminProfile | null;
  loading: boolean;
  isSuperadmin: boolean;
  refresh: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = React.createContext<AuthState | null>(null);

/**
 * Provides the current admin profile (from /admin/api/me) to the panel. The
 * middleware already blocks unauthenticated access to /admin pages; this adds
 * the "Signed in as" identity and the superadmin gate for the UI.
 */
export function AdminAuthProvider({ children }: { children: React.ReactNode }) {
  const [admin, setAdmin] = React.useState<AdminProfile | null>(null);
  const [loading, setLoading] = React.useState(true);
  const router = useRouter();

  const refresh = React.useCallback(async () => {
    try {
      const res = await fetch("/admin/api/me", { cache: "no-store" });
      if (!res.ok) {
        setAdmin(null);
        router.replace("/admin/login");
        return;
      }
      setAdmin((await res.json()) as AdminProfile);
    } catch {
      setAdmin(null);
    } finally {
      setLoading(false);
    }
  }, [router]);

  const logout = React.useCallback(async () => {
    await fetch("/admin/api/logout", { method: "POST" });
    setAdmin(null);
    router.replace("/admin/login");
  }, [router]);

  React.useEffect(() => {
    void refresh();
  }, [refresh]);

  const value: AuthState = {
    admin,
    loading,
    isSuperadmin: admin?.role === "superadmin",
    refresh,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAdminAuth(): AuthState {
  const ctx = React.useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAdminAuth must be used within AdminAuthProvider");
  }
  return ctx;
}
