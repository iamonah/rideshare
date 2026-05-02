import type { Metadata } from "next";
import Script from "next/script";
import "./globals.css";

export const metadata: Metadata = {
  title: "RideShare",
  description: "RideShare",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link
          rel="stylesheet"
          href="https://unpkg.com/maplibre-gl/dist/maplibre-gl.css"
        />
      </head>
      <body className="antialiased">
        <Script
          src="https://unpkg.com/maplibre-gl/dist/maplibre-gl.js"
          strategy="beforeInteractive"
        />
        {children}
      </body>
    </html>
  );
}
