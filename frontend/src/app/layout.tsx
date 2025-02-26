"use client";

import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AppBar, Toolbar, Typography, Button, Chip } from "@mui/material";
import toast, { Toaster } from "react-hot-toast";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";
const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [translating, setTranslating] = useState("");
  const [toastId, setToastId] = useState<any>();
  const toastIdRef = useRef<string>("");

  useEffect(() => {
    const webSocket = new WebSocket("ws://localhost:8080/ws");
    webSocket.onmessage = function (message) {
        console.log(message)
      switch (message.data.slice(0, 1)) {
        case '1':
          toast.loading(message.data.slice(1), {
            id: toastIdRef.current,
          });
        case '2':
        case '3':
          toast.success(message.data.slice(1), {
            id: toastIdRef.current,
          });
      }
    };
    return () => webSocket.close();
  }, []);

  const handleTranslate = async () => {
    if (!toastId) {
      const toastId = toast.loading("Translating...");
      toastIdRef.current = toastId;
    } else {
      toast.loading("Translating...", {
        id: toastIdRef.current,
      });
    }
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/translate`);
  };

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div>
          <Toaster />
        </div>
        <AppBar position="static">
          <Toolbar>
            <Typography variant="h6" sx={{ flexGrow: 1 }}>
              Melons Quiz App
            </Typography>
            <Typography variant="h6" sx={{ flexGrow: 1 }}>
              {translating}
            </Typography>
            <Button color="inherit" onClick={handleTranslate}>
              Translate
            </Button>
            <Button color="inherit" component={Link} href="/home">
              Home
            </Button>
            <Button color="inherit" component={Link} href="/quiz">
              Quiz
            </Button>
            <Button color="inherit" component={Link} href="/manage">
              Manage
            </Button>
          </Toolbar>
        </AppBar>
        <div>{children}</div>
      </body>
    </html>
  );
}
