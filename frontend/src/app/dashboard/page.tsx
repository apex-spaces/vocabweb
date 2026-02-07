'use client'

import { useEffect, useState } from 'react'
import Header from "@/components/layout/Header"
import Sidebar from "@/components/layout/Sidebar"
import StatsCard from "@/components/dashboard/StatsCard"
import WeeklyChart from "@/components/dashboard/WeeklyChart"

interface DashboardData {
  today_due: number
  today_new: number
  total_mastered: number
  streak_days: number
  recent_words: Array<{
    word_id: number
    word: string
    definition: string
    created_at: string
  }>
  weekly_stats: Array<{
    date: string
    count: number
  }>
}

export default function DashboardPage() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    try {
      const response = await fetch('/api/v1/dashboard', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      
      if (response.ok) {
        const result = await response.json()
        setData(result)
      }
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex min-h-screen bg-[#0F172A]">
        <Sidebar />
        <div className="flex-1">
          <Header />
          <main className="p-8">
            <div className="text-gray-400">Loading...</div>
          </main>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen bg-[#0F172A]">
      <Sidebar />
      <div className="flex-1">
        <Header />
        <main className="p-8">
          <h1 className="text-4xl font-serif font-bold mb-8 text-white" style={{ fontFamily: 'Playfair Display, serif' }}>
            Dashboard
          </h1>

          {/* Stats Cards Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <StatsCard 
              title="Due Today" 
              value={data?.today_due || 0} 
              icon="📚" 
              color="#F59E0B" 
            />
            <StatsCard 
              title="New Today" 
              value={data?.today_new || 0} 
              icon="✨" 
              color="#10B981" 
            />
            <StatsCard 
              title="Mastered" 
              value={data?.total_mastered || 0} 
              icon="🎯" 
              color="#8B5CF6" 
            />
            <StatsCard 
              title="Streak Days" 
              value={data?.streak_days || 0} 
              icon="🔥" 
              color="#EF4444" 
            />
          </div>

          {/* Start Review Button */}
          {data && data.today_due > 0 && (
            <div className="mb-8">
              <button className="relative bg-[#F59E0B] hover:bg-[#FBBF24] text-white font-semibold px-8 py-4 rounded-lg transition-all duration-300 shadow-lg hover:shadow-xl">
                <span className="relative z-10">Start Review ({data.today_due} words)</span>
                <span className="absolute inset-0 rounded-lg bg-[#F59E0B] animate-pulse opacity-75"></span>
              </button>
            </div>
          )}

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Recent Words */}
            <div className="bg-[#1E293B] rounded-lg p-6 border border-gray-800">
              <h3 className="text-xl font-semibold text-white mb-4">Recent Words</h3>
              <div className="space-y-3">
                {data?.recent_words && data.recent_words.length > 0 ? (
                  data.recent_words.map((word) => (
                    <div key={word.word_id} className="p-3 bg-[#0F172A] rounded border border-gray-800 hover:border-gray-700 transition-colors">
                      <div className="font-semibold text-white mb-1">{word.word}</div>
                      <div className="text-sm text-gray-400">{word.definition}</div>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-400">No words collected yet</p>
                )}
              </div>
            </div>

            {/* Weekly Chart */}
            <WeeklyChart data={data?.weekly_stats || []} />
          </div>
        </main>
      </div>
    </div>
  )
}
