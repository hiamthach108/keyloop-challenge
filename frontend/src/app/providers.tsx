'use client';

import { AntdRegistry } from '@ant-design/nextjs-registry';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { App, ConfigProvider } from 'antd';
import { useState, type ReactNode } from 'react';
import { Toaster } from 'react-hot-toast';

import { TOAST_DEFAULT_OPTIONS } from '@/config/helpers/toast.helper';

export default function AppProvider({ children }: Readonly<{ children: ReactNode }>) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        staleTime: 30_000,
      },
    },
  }));

  return (
    <AntdRegistry>
      <QueryClientProvider client={queryClient}>
        <ConfigProvider
          theme={{
            token: {
              borderRadius: 8,
              colorPrimary: '#1677ff',
            },
          }}
        >
          <App>
            {children}
            <Toaster {...TOAST_DEFAULT_OPTIONS} />
          </App>
        </ConfigProvider>
      </QueryClientProvider>
    </AntdRegistry>
  );
}
