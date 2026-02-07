interface StatsCardProps {
  title: string
  value: number
  icon: string
  color: string
}

export default function StatsCard({ title, value, icon, color }: StatsCardProps) {
  return (
    <div className="bg-[#1E293B] rounded-lg p-6 border border-gray-800 hover:border-gray-700 transition-colors">
      <div className="flex items-start justify-between">
        <div>
          <p className="text-gray-400 text-sm mb-2">{title}</p>
          <p className="text-3xl font-bold text-white">{value}</p>
        </div>
        <div 
          className="text-3xl"
          style={{ color }}
        >
          {icon}
        </div>
      </div>
    </div>
  )
}
