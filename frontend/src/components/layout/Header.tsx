import Logo from "@/components/common/Logo"

export default function Header() {
  return (
    <header className="bg-card border-b border-gray-700 px-6 py-4">
      <div className="flex items-center justify-between">
        <Logo />
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-400">Welcome back!</span>
        </div>
      </div>
    </header>
  )
}
