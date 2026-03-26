import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ 
  subsets: ["latin"],
  variable: '--font-inter',
});

export const metadata = {
  title: "gitresolve",
  description: "A secure, offline deterministic Git conflict resolution engine.",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={inter.variable}>
      <body className="bg-black text-[#ededed] antialiased">
        {children}
      </body>
    </html>
  );
}