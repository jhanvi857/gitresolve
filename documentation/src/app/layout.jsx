import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ 
  subsets: ["latin"],
  display: 'swap',
});

export const metadata = {
  title: "gitresolve – Deterministic Git Conflict Resolution",
  description: "A purely offline, deterministic AST-powered Git conflict resolution engine.",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={inter.className}>
      <body className="bg-black text-[#ededed] antialiased selection:bg-blue-500/30 selection:text-white">
        {children}
      </body>
    </html>
  );
}