import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import Providers from '@/components/layout/Providers'
import './globals.css'

const inter = Inter({ subsets: ['latin', 'vietnamese'] })

export const metadata: Metadata = {
  title: 'UTM - Hệ thống Quản lý Không phận UAV',
  description: 'Unmanned Traffic Management System',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="vi">
      <body className={inter.className}>
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
