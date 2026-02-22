/**
 * GoogleOAuthWrapper — Client component that wraps the app with GoogleOAuthProvider.
 *
 * Separated from layout.tsx because GoogleOAuthProvider requires "use client".
 */
"use client";

import { GoogleOAuthProvider } from "@react-oauth/google";

const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "";

export function GoogleOAuthWrapper({ children }: { children: React.ReactNode }) {
  if (!GOOGLE_CLIENT_ID) {
    // If no Google Client ID is configured, render children without the provider
    return <>{children}</>;
  }

  return (
    <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
      {children}
    </GoogleOAuthProvider>
  );
}
