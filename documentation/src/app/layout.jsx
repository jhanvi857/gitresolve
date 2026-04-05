import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ 
  subsets: ["latin"],
  variable: '--font-inter',
});

export const metadata = {
  title: "gitresolve",
  description: "A secure, offline deterministic Git conflict resolution engine.",
  icons: {
    icon: "/logo.png",
  },
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={inter.variable}>
      <head>
        <link rel="icon" href="/favicon.ico" sizes="any" />
        <link rel="icon" href="/logo.png" type="image/png" />
      </head>
      <body className="bg-black text-[#ededed] antialiased">
        {children}
      </body>
    </html>
  );
}