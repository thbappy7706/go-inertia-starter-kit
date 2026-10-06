import React from "react"
import { Link, usePage } from "@inertiajs/react"
import { LayoutDashboard, Package, Users, Layers } from "lucide-react"

interface NavItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}

const navItems: NavItem[] = [
  { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { name: "Products", href: "/products", icon: Package },
  { name: "Users", href: "/users", icon: Users },
]

interface SidebarProps {
  onItemClick?: () => void
}

export const Sidebar: React.FC<SidebarProps> = ({ onItemClick }) => {
  const { url } = usePage()

  return (
    <aside className="w-64 flex flex-col h-full bg-slate-900 text-slate-100 border-r border-slate-800">
      {/* Brand logo */}
      <div className="h-16 flex items-center gap-3 px-6 border-b border-slate-800">
        <div className="h-9 w-9 rounded-lg bg-indigo-600 flex items-center justify-center text-white shadow-md">
          <Layers className="h-5 w-5" />
        </div>
        <div>
          <span className="font-bold text-base tracking-tight text-white block">Starter Kit</span>
          <span className="text-[10px] text-slate-400 font-medium tracking-wider uppercase">Go + Inertia v3</span>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-3 py-6 space-y-1.5 overflow-y-auto">
        <div className="px-3 pb-2 text-[11px] font-semibold tracking-wider text-slate-400 uppercase">
          Management
        </div>

        {navItems.map((item) => {
          const Icon = item.icon
          const isActive = url === item.href || (item.href !== "/dashboard" && url.startsWith(item.href))

          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={onItemClick}
              className={`flex items-center gap-3 px-3.5 py-2.5 rounded-lg text-sm font-medium transition-all ${
                isActive
                  ? "bg-indigo-600 text-white shadow-sm"
                  : "text-slate-300 hover:bg-slate-800/80 hover:text-white"
              }`}
            >
              <Icon className={`h-4 w-4 ${isActive ? "text-white" : "text-slate-400"}`} />
              <span>{item.name}</span>
            </Link>
          )
        })}
      </nav>

      {/* Footer footer info */}
      <div className="p-4 border-t border-slate-800">
        <div className="rounded-lg bg-slate-800/60 p-3 text-xs text-slate-400 flex items-center justify-between">
          <div>
            <div className="font-semibold text-slate-200">PostgreSQL 17</div>
            <div className="text-[11px] text-emerald-400 font-medium">Connected</div>
          </div>
          <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
        </div>
      </div>
    </aside>
  )
}