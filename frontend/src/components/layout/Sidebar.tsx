import Link from "next/link"
import { Home, BookOpen, Brain, Settings, BarChart, Users, FileText, Calendar, Trophy, Target, Bell, HelpCircle } from "lucide-react"

const navItems = [
  { href: "/dashboard", icon: Home, label: "Dashboard" },
  { href: "/words", icon: BookOpen, label: "Words" },
  { href: "/review", icon: Brain, label: "Review" },
  { href: "/progress", icon: BarChart, label: "Progress" },
  { href: "/collections", icon: FileText, label: "Collections" },
  { href: "/schedule", icon: Calendar, label: "Schedule" },
  { href: "/achievements", icon: Trophy, label: "Achievements" },
  { href: "/goals", icon: Target, label: "Goals" },
  { href: "/community", icon: Users, label: "Community" },
  { href: "/notifications", icon: Bell, label: "Notifications" },
  { href: "/help", icon: HelpCircle, label: "Help" },
  { href: "/settings", icon: Settings, label: "Settings" },
]

export default function Sidebar() {
  return (
    <aside className="w-64 bg-card border-r border-gray-700 min-h-screen p-4">
      <nav className="space-y-2">
        {navItems.map((item) => {
          const Icon = item.icon
          return (
            <Link
              key={item.href}
              href={item.href}
              className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-background transition-colors text-gray-300 hover:text-accent"
            >
              <Icon className="w-5 h-5" />
              <span>{item.label}</span>
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
