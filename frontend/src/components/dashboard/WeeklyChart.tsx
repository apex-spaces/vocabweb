interface WeeklyChartProps {
  data: Array<{ date: string; count: number }>
}

export default function WeeklyChart({ data }: WeeklyChartProps) {
  const maxCount = Math.max(...data.map(d => d.count), 1)
  
  return (
    <div className="bg-[#1E293B] rounded-lg p-6 border border-gray-800">
      <h3 className="text-xl font-semibold text-white mb-6">Weekly Activity</h3>
      
      <div className="flex items-end justify-between gap-3 h-48">
        {data.map((day, index) => {
          const height = (day.count / maxCount) * 100
          
          return (
            <div key={index} className="flex-1 flex flex-col items-center gap-2">
              <div className="relative w-full flex items-end justify-center h-40">
                <div 
                  className="w-full bg-[#F59E0B] rounded-t transition-all duration-300 hover:bg-[#FBBF24]"
                  style={{ height: `${height}%` }}
                  title={`${day.count} words`}
                />
              </div>
              <span className="text-xs text-gray-400">{day.date}</span>
            </div>
          )
        })}
      </div>
    </div>
  )
}
