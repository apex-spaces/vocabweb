'use client'

import { useState, useEffect, useCallback } from 'react'
import Header from "@/components/layout/Header"
import Sidebar from "@/components/layout/Sidebar"

interface Word {
  user_word_id: string
  word_id: string
  word: string
  phonetic: string
  definitions: string
  easiness_factor: number
  interval: number
  repetitions: number
  context_sentence: string
}

interface ReviewStats {
  total_due: number
  reviewed: number
  new_words: number
  mastered_today: number
}

export default function ReviewPage() {
  const [words, setWords] = useState<Word[]>([])
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isFlipped, setIsFlipped] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [stats, setStats] = useState<ReviewStats | null>(null)
  const [isComplete, setIsComplete] = useState(false)

  // Fetch due words on mount
  useEffect(() => {
    fetchDueWords()
    fetchStats()
  }, [])

  const fetchDueWords = async () => {
    try {
      setIsLoading(true)
      const response = await fetch('/api/v1/review/due?limit=20', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      const data = await response.json()
      if (data.success) {
        setWords(data.data.words || [])
        if (data.data.words.length === 0) {
          setIsComplete(true)
        }
      }
    } catch (error) {
      console.error('Failed to fetch due words:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const fetchStats = async () => {
    try {
      const response = await fetch('/api/v1/review/stats', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      const data = await response.json()
      if (data.success) {
        setStats(data.data)
      }
    } catch (error) {
      console.error('Failed to fetch stats:', error)
    }
  }

  const submitReview = async (quality: number) => {
    if (isSubmitting || currentIndex >= words.length) return

    const currentWord = words[currentIndex]
    setIsSubmitting(true)

    try {
      const response = await fetch('/api/v1/review/submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          user_word_id: currentWord.user_word_id,
          quality: quality
        })
      })

      if (response.ok) {
        // Move to next word
        if (currentIndex + 1 >= words.length) {
          setIsComplete(true)
          fetchStats()
        } else {
          setCurrentIndex(currentIndex + 1)
          setIsFlipped(false)
        }
      }
    } catch (error) {
      console.error('Failed to submit review:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleFlip = () => {
    setIsFlipped(!isFlipped)
  }

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      if (isSubmitting || isComplete) return

      if (e.code === 'Space') {
        e.preventDefault()
        handleFlip()
      } else if (isFlipped) {
        if (e.key === '1') submitReview(1) // Don't know
        if (e.key === '2') submitReview(3) // Vague
        if (e.key === '3') submitReview(5) // Know
      }
    }

    window.addEventListener('keydown', handleKeyPress)
    return () => window.removeEventListener('keydown', handleKeyPress)
  }, [isFlipped, isSubmitting, isComplete, currentIndex])

  if (isLoading) {
    return (
      <div className="flex min-h-screen bg-gray-900">
        <Sidebar />
        <div className="flex-1">
          <Header />
          <main className="p-8 flex items-center justify-center">
            <div className="text-gray-400 text-xl">Loading...</div>
          </main>
        </div>
      </div>
    )
  }

  if (isComplete || words.length === 0) {
    return (
      <div className="flex min-h-screen bg-gray-900">
        <Sidebar />
        <div className="flex-1">
          <Header />
          <main className="p-8">
            <div className="max-w-2xl mx-auto text-center">
              <div className="bg-gray-800 rounded-2xl p-12 border border-gray-700">
                <div className="text-6xl mb-6">🎉</div>
                <h1 className="text-4xl font-bold text-white mb-4">
                  {words.length === 0 ? 'No Words to Review' : 'Review Complete!'}
                </h1>
                <p className="text-gray-400 text-lg mb-8">
                  {words.length === 0 
                    ? 'You have no words due for review right now. Come back later!'
                    : 'Great job! You\'ve completed all your reviews for now.'}
                </p>
                
                {stats && (
                  <div className="grid grid-cols-2 gap-4 mb-8">
                    <div className="bg-gray-700/50 rounded-lg p-4">
                      <div className="text-3xl font-bold text-blue-400">{stats.reviewed}</div>
                      <div className="text-sm text-gray-400">Reviewed Today</div>
                    </div>
                    <div className="bg-gray-700/50 rounded-lg p-4">
                      <div className="text-3xl font-bold text-green-400">{stats.mastered_today}</div>
                      <div className="text-sm text-gray-400">Mastered Today</div>
                    </div>
                  </div>
                )}

                <button
                  onClick={() => window.location.href = '/dashboard'}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-lg font-medium transition-colors"
                >
                  Back to Dashboard
                </button>
              </div>
            </div>
          </main>
        </div>
      </div>
    )
  }

  const currentWord = words[currentIndex]
  const progress = ((currentIndex + 1) / words.length) * 100

  return (
    <div className="flex min-h-screen bg-gray-900">
      <Sidebar />
      <div className="flex-1">
        <Header />
        <main className="p-8">
          <div className="max-w-4xl mx-auto">
            {/* Progress Bar */}
            <div className="mb-8">
              <div className="flex justify-between text-sm text-gray-400 mb-2">
                <span>Progress</span>
                <span>{currentIndex + 1} / {words.length}</span>
              </div>
              <div className="w-full bg-gray-800 rounded-full h-3 overflow-hidden">
                <div 
                  className="bg-blue-600 h-full transition-all duration-300 ease-out"
                  style={{ width: `${progress}%` }}
                />
              </div>
            </div>

            {/* Review Card */}
            <div className="relative" style={{ perspective: '1000px' }}>
              <div 
                className={`relative w-full transition-transform duration-500 ${
                  isFlipped ? 'rotate-y-180' : ''
                }`}
                style={{ transformStyle: 'preserve-3d' }}
                onClick={handleFlip}
              >
                {/* Front Side - Word */}
                <div 
                  className="w-full"
                  style={{ 
                    backfaceVisibility: 'hidden',
                    display: isFlipped ? 'none' : 'block'
                  }}
                >
                  <div className="bg-gray-800 rounded-2xl p-12 border border-gray-700 min-h-[400px] flex flex-col items-center justify-center cursor-pointer hover:border-gray-600 transition-colors">
                    <div className="text-center">
                      <h2 className="text-6xl font-bold text-white mb-4">
                        {currentWord.word}
                      </h2>
                      {currentWord.phonetic && (
                        <p className="text-xl text-gray-400 mb-6">
                          /{currentWord.phonetic}/
                        </p>
                      )}
                      {currentWord.context_sentence && (
                        <div className="mt-8 p-4 bg-gray-700/50 rounded-lg max-w-2xl">
                          <p className="text-sm text-gray-400 mb-1">Context:</p>
                          <p className="text-gray-300 italic">"{currentWord.context_sentence}"</p>
                        </div>
                      )}
                    </div>
                    <div className="mt-8 text-gray-500 text-sm">
                      Click or press <kbd className="px-2 py-1 bg-gray-700 rounded">Space</kbd> to flip
                    </div>
                  </div>
                </div>

                {/* Back Side - Definition */}
                <div 
                  className="w-full"
                  style={{ 
                    backfaceVisibility: 'hidden',
                    display: isFlipped ? 'block' : 'none'
                  }}
                >
                  <div className="bg-gray-800 rounded-2xl p-12 border border-gray-700 min-h-[400px]">
                    <div className="text-center mb-8">
                      <h2 className="text-4xl font-bold text-white mb-4">
                        {currentWord.word}
                      </h2>
                      {currentWord.phonetic && (
                        <p className="text-lg text-gray-400 mb-6">
                          /{currentWord.phonetic}/
                        </p>
                      )}
                    </div>

                    {/* Definitions */}
                    <div className="mb-8 max-w-2xl mx-auto">
                      <div className="bg-gray-700/50 rounded-lg p-6">
                        <div 
                          className="text-gray-300 space-y-3"
                          dangerouslySetInnerHTML={{ 
                            __html: formatDefinitions(currentWord.definitions) 
                          }}
                        />
                      </div>
                    </div>

                    <div className="text-center text-gray-500 text-sm mb-6">
                      How well did you remember this word?
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Rating Buttons (only show when flipped) */}
            {isFlipped && (
              <div className="mt-8 grid grid-cols-3 gap-4">
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    submitReview(1)
                  }}
                  disabled={isSubmitting}
                  className="bg-red-600 hover:bg-red-700 disabled:bg-red-800 disabled:cursor-not-allowed text-white py-4 px-6 rounded-xl font-medium transition-colors flex flex-col items-center gap-2"
                >
                  <span className="text-2xl">😕</span>
                  <span>Don't Know</span>
                  <kbd className="text-xs bg-red-700 px-2 py-1 rounded">1</kbd>
                </button>

                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    submitReview(3)
                  }}
                  disabled={isSubmitting}
                  className="bg-yellow-600 hover:bg-yellow-700 disabled:bg-yellow-800 disabled:cursor-not-allowed text-white py-4 px-6 rounded-xl font-medium transition-colors flex flex-col items-center gap-2"
                >
                  <span className="text-2xl">🤔</span>
                  <span>Vague</span>
                  <kbd className="text-xs bg-yellow-700 px-2 py-1 rounded">2</kbd>
                </button>

                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    submitReview(5)
                  }}
                  disabled={isSubmitting}
                  className="bg-green-600 hover:bg-green-700 disabled:bg-green-800 disabled:cursor-not-allowed text-white py-4 px-6 rounded-xl font-medium transition-colors flex flex-col items-center gap-2"
                >
                  <span className="text-2xl">✅</span>
                  <span>Know Well</span>
                  <kbd className="text-xs bg-green-700 px-2 py-1 rounded">3</kbd>
                </button>
              </div>
            )}

            {/* Keyboard Shortcuts Hint */}
            <div className="mt-6 text-center text-gray-500 text-sm">
              <p>Keyboard shortcuts: <kbd className="px-2 py-1 bg-gray-800 rounded">Space</kbd> to flip, <kbd className="px-2 py-1 bg-gray-800 rounded">1</kbd>/<kbd className="px-2 py-1 bg-gray-800 rounded">2</kbd>/<kbd className="px-2 py-1 bg-gray-800 rounded">3</kbd> to rate</p>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}

// Helper function to format definitions JSON
function formatDefinitions(definitionsJson: string): string {
  try {
    const definitions = JSON.parse(definitionsJson)
    if (Array.isArray(definitions)) {
      return definitions.map((def: any, index: number) => {
        return `<div class="mb-3">
          <span class="font-semibold text-blue-400">${def.pos || 'n.'}</span>
          <p class="mt-1">${def.meaning || def.definition || ''}</p>
          ${def.example ? `<p class="mt-1 text-sm text-gray-400 italic">"${def.example}"</p>` : ''}
        </div>`
      }).join('')
    }
    return definitionsJson
  } catch {
    return definitionsJson
  }
}
