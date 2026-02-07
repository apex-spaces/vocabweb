import Header from "@/components/layout/Header"
import Sidebar from "@/components/layout/Sidebar"

export default function HelpPage() {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1">
        <Header />
        <main className="p-8">
          <h1 className="text-4xl font-serif font-bold mb-4">Help</h1>
          <div className="text-gray-400">
            <p className="text-lg">Coming Soon</p>
            <p className="mt-2">Help documentation and FAQs will appear here.</p>
          </div>
        </main>
      </div>
    </div>
  )
}
