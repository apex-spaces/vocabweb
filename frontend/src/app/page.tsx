import Link from "next/link"
import Logo from "@/components/common/Logo"

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="p-6 border-b border-gray-700">
        <Logo />
      </header>
      
      <main className="flex-1 flex items-center justify-center px-6">
        <div className="max-w-2xl text-center space-y-8">
          <h1 className="text-5xl font-serif font-bold text-accent">
            Master Your Vocabulary
          </h1>
          <p className="text-xl text-gray-300">
            Build your word power with intelligent spaced repetition and personalized learning paths.
          </p>
          <div className="flex gap-4 justify-center">
            <Link
              href="/auth"
              className="px-8 py-3 bg-accent text-background font-semibold rounded-lg hover:bg-amber-500 transition-colors"
            >
              Get Started
            </Link>
            <Link
              href="/dashboard"
              className="px-8 py-3 bg-card text-gray-100 font-semibold rounded-lg hover:bg-gray-700 transition-colors"
            >
              View Demo
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
