'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/lib/store/auth'
import { clsx } from 'clsx'

const navItems = [
  { href: '/dashboard', label: 'Tổng quan', roles: ['admin', 'pilot', 'operator', 'atc', 'inspector'] },
  { href: '/monitoring', label: 'Giám sát bay', roles: ['admin', 'atc', 'inspector'] },
  { href: '/aircraft', label: 'Máy bay', roles: ['admin', 'pilot', 'operator', 'atc', 'inspector'] },
  { href: '/authorization', label: 'Cấp phép bay', roles: ['admin', 'pilot', 'operator', 'atc', 'inspector'] },
  { href: '/airspace', label: 'Vùng không phận', roles: ['admin', 'atc', 'inspector'] },
  { href: '/notifications', label: 'Thông báo', roles: ['admin', 'pilot', 'operator', 'atc', 'inspector'] },
] as const

export default function Sidebar() {
  const pathname = usePathname()
  const { user, logout, hasRole } = useAuthStore()

  return (
    <aside className="flex flex-col w-64 min-h-screen bg-gray-900 text-white">
      <div className="p-6 border-b border-gray-700">
        <h1 className="text-xl font-bold text-blue-400">UTM System</h1>
        <p className="text-xs text-gray-400 mt-1">Quản lý Không phận UAV</p>
      </div>

      <nav className="flex-1 p-4 space-y-1">
        {navItems
          .filter((item) => hasRole(...(item.roles as any[])))
          .map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={clsx(
                'block px-4 py-2.5 rounded-lg text-sm transition-colors',
                pathname.startsWith(item.href)
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-300 hover:bg-gray-800'
              )}
            >
              {item.label}
            </Link>
          ))}
      </nav>

      <div className="p-4 border-t border-gray-700">
        <div className="text-sm text-gray-300 mb-2 truncate">{user?.full_name}</div>
        <div className="text-xs text-gray-500 mb-3 uppercase">{user?.role}</div>
        <button
          onClick={logout}
          className="w-full px-4 py-2 text-sm text-red-400 border border-red-400 rounded-lg hover:bg-red-400 hover:text-white transition-colors"
        >
          Đăng xuất
        </button>
      </div>
    </aside>
  )
}
