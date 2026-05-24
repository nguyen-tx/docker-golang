import { Metadata } from 'next'

export const metadata: Metadata = { title: 'Tổng quan - UTM' }

export default function DashboardPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Tổng quan hệ thống</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard label="Máy bay đăng ký" value="--" color="blue" />
        <StatCard label="Đơn chờ xét duyệt" value="--" color="yellow" />
        <StatCard label="Chuyến bay đang hoạt động" value="--" color="green" />
        <StatCard label="Vi phạm hôm nay" value="--" color="red" />
      </div>
    </div>
  )
}

function StatCard({ label, value, color }: { label: string; value: string; color: string }) {
  const colors: Record<string, string> = {
    blue: 'bg-blue-50 border-blue-200 text-blue-700',
    yellow: 'bg-yellow-50 border-yellow-200 text-yellow-700',
    green: 'bg-green-50 border-green-200 text-green-700',
    red: 'bg-red-50 border-red-200 text-red-700',
  }
  return (
    <div className={`p-6 rounded-xl border-2 ${colors[color]}`}>
      <div className="text-3xl font-bold">{value}</div>
      <div className="text-sm mt-1 opacity-80">{label}</div>
    </div>
  )
}
