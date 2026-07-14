import type { Metadata } from 'next';
import type { ReactNode } from 'react';

import MainLayout from '@/components/layout/MainLayout';

import AppProvider from './providers';

import '@/styles/globals.css';

export const metadata: Metadata = {
  title: 'Inventory Dashboard | Keyloop Challenge',
  description: 'Intelligent dealership inventory dashboard for aging stock.',
};

export default function RootLayout({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="en">
      <body suppressHydrationWarning>
        <AppProvider>
          <MainLayout>{children}</MainLayout>
        </AppProvider>
      </body>
    </html>
  );
}
