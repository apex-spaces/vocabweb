export interface User {
  id: string
  email: string
  name: string
  createdAt: string
}

export interface Word {
  id: string
  word: string
  definition: string
  example?: string
  createdAt: string
}

export interface ReviewSession {
  id: string
  userId: string
  wordId: string
  result: 'correct' | 'incorrect'
  createdAt: string
}
