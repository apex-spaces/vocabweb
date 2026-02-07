'use client'

import { useState, useEffect } from 'react'
import Header from "@/components/layout/Header"
import Sidebar from "@/components/layout/Sidebar"

interface Word {
  id: number
  word: string
  definition?: string
  source?: string
  created_at: string
  status: string
}

interface WordCandidate {
  word: string
  frequency: number
  is_collected: boolean
  word_id?: number
}

export default function WordsPage() {
  const [words, setWords] = useState<Word[]>([])
  const [loading, setLoading] = useState(false)
  const [showAddModal, setShowAddModal] = useState(false)
  const [showAnalyzeModal, setShowAnalyzeModal] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [limit] = useState(20)

  // Add word form state
  const [newWord, setNewWord] = useState({
    word: '',
    definition: '',
    source: 'manual',
    context: ''
  })

  // Analyze text state
  const [analyzeText, setAnalyzeText] = useState('')
  const [candidates, setCandidates] = useState<WordCandidate[]>([])
  const [analyzing, setAnalyzing] = useState(false)

  // Fetch words
  useEffect(() => {
    fetchWords()
  }, [page, searchQuery])

  const fetchWords = async () => {
    setLoading(true)
    try {
      const token = localStorage.getItem('token')
      const response = await fetch(
        `http://localhost:8080/api/v1/words?page=${page}&limit=${limit}&sort=created_at&order=desc`,
        {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      )
      const data = await response.json()
      setWords(data.words || [])
    } catch (error) {
      console.error('Failed to fetch words:', error)
    } finally {
      setLoading(false)
    }
  }

  // Add single word
  const handleAddWord = async () => {
    if (!newWord.word) return

    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/words', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(newWord)
      })

      if (response.ok) {
        setShowAddModal(false)
        setNewWord({ word: '', definition: '', source: 'manual', context: '' })
        fetchWords()
      }
    } catch (error) {
      console.error('Failed to add word:', error)
    }
  }

  // Analyze text
  const handleAnalyzeText = async () => {
    if (!analyzeText) return

    setAnalyzing(true)
    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/words/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ text: analyzeText })
      })

      const data = await response.json()
      setCandidates(data.candidates || [])
    } catch (error) {
      console.error('Failed to analyze text:', error)
    } finally {
      setAnalyzing(false)
    }
  }

  // Delete word
  const handleDeleteWord = async (wordId: number) => {
    if (!confirm('Are you sure you want to delete this word?')) return

    try {
      const token = localStorage.getItem('token')
      const response = await fetch(`http://localhost:8080/api/v1/words/${wordId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (response.ok) {
        fetchWords()
      }
    } catch (error) {
      console.error('Failed to delete word:', error)
    }
  }

  return (
    <div className="flex min-h-screen bg-gray-900">
      <Sidebar />
      <div className="flex-1">
        <Header />
        <main className="p-8">
          {/* Page Header */}
          <div className="mb-8">
            <h1 className="text-4xl font-serif font-bold text-white mb-2">My Words</h1>
            <p className="text-gray-400">Manage your vocabulary collection</p>
          </div>

          {/* Action Bar */}
          <div className="mb-6 flex gap-4">
            <input
              type="text"
              placeholder="Search words..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="flex-1 px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
            />
            <button
              onClick={() => setShowAddModal(true)}
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
            >
              + Add Word
            </button>
            <button
              onClick={() => setShowAnalyzeModal(true)}
              className="px-6 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium transition-colors"
            >
              📝 Analyze Text
            </button>
          </div>

          {/* Words List */}
          {loading ? (
            <div className="text-center text-gray-400 py-12">Loading...</div>
          ) : words.length === 0 ? (
            <div className="text-center text-gray-400 py-12">
              <p className="text-lg">No words yet</p>
              <p className="mt-2">Add your first word or analyze some text to get started</p>
            </div>
          ) : (
            <div className="bg-gray-800 rounded-lg overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-700">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Word</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Definition</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Source</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Added</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-700">
                  {words.map((word) => (
                    <tr key={word.id} className="hover:bg-gray-750">
                      <td className="px-6 py-4 whitespace-nowrap text-white font-medium">{word.word}</td>
                      <td className="px-6 py-4 text-gray-300">{word.definition || '-'}</td>
                      <td className="px-6 py-4 text-gray-400 text-sm">{word.source || '-'}</td>
                      <td className="px-6 py-4">
                        <span className="px-2 py-1 text-xs rounded-full bg-blue-900 text-blue-300">
                          {word.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-gray-400 text-sm">
                        {new Date(word.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4">
                        <button
                          onClick={() => handleDeleteWord(word.id)}
                          className="text-red-400 hover:text-red-300 text-sm"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Pagination */}
          <div className="mt-6 flex justify-center gap-2">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-4 py-2 bg-gray-800 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-700"
            >
              Previous
            </button>
            <span className="px-4 py-2 text-gray-400">Page {page}</span>
            <button
              onClick={() => setPage(p => p + 1)}
              disabled={words.length < limit}
              className="px-4 py-2 bg-gray-800 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-700"
            >
              Next
            </button>
          </div>

          {/* Add Word Modal */}
          {showAddModal && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
              <div className="bg-gray-800 rounded-lg p-6 w-full max-w-md">
                <h2 className="text-2xl font-bold text-white mb-4">Add New Word</h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-gray-300 mb-2">Word *</label>
                    <input
                      type="text"
                      value={newWord.word}
                      onChange={(e) => setNewWord({ ...newWord, word: e.target.value })}
                      className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                      placeholder="Enter word"
                    />
                  </div>
                  <div>
                    <label className="block text-gray-300 mb-2">Definition</label>
                    <textarea
                      value={newWord.definition}
                      onChange={(e) => setNewWord({ ...newWord, definition: e.target.value })}
                      className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                      placeholder="Enter definition"
                      rows={3}
                    />
                  </div>
                  <div>
                    <label className="block text-gray-300 mb-2">Context</label>
                    <input
                      type="text"
                      value={newWord.context}
                      onChange={(e) => setNewWord({ ...newWord, context: e.target.value })}
                      className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                      placeholder="Where did you see this word?"
                    />
                  </div>
                  <div className="flex gap-3 mt-6">
                    <button
                      onClick={handleAddWord}
                      className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium"
                    >
                      Add Word
                    </button>
                    <button
                      onClick={() => setShowAddModal(false)}
                      className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Analyze Text Modal */}
          {showAnalyzeModal && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
              <div className="bg-gray-800 rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-y-auto">
                <h2 className="text-2xl font-bold text-white mb-4">Analyze Text</h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-gray-300 mb-2">Paste English text</label>
                    <textarea
                      value={analyzeText}
                      onChange={(e) => setAnalyzeText(e.target.value)}
                      className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
                      placeholder="Paste your English text here..."
                      rows={6}
                    />
                  </div>
                  <button
                    onClick={handleAnalyzeText}
                    disabled={analyzing || !analyzeText}
                    className="w-full px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {analyzing ? 'Analyzing...' : 'Analyze'}
                  </button>

                  {/* Candidates List */}
                  {candidates.length > 0 && (
                    <div className="mt-6">
                      <h3 className="text-lg font-semibold text-white mb-3">Found {candidates.length} words</h3>
                      <div className="space-y-2 max-h-64 overflow-y-auto">
                        {candidates.map((candidate, index) => (
                          <div
                            key={index}
                            className="flex items-center justify-between p-3 bg-gray-700 rounded-lg"
                          >
                            <div className="flex-1">
                              <span className="text-white font-medium">{candidate.word}</span>
                              <span className="ml-3 text-gray-400 text-sm">
                                appears {candidate.frequency}x
                              </span>
                            </div>
                            {candidate.is_collected ? (
                              <span className="text-green-400 text-sm">✓ Collected</span>
                            ) : (
                              <span className="text-gray-500 text-sm">New</span>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  <div className="flex gap-3 mt-6">
                    <button
                      onClick={() => {
                        setShowAnalyzeModal(false)
                        setAnalyzeText('')
                        setCandidates([])
                      }}
                      className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium"
                    >
                      Close
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
