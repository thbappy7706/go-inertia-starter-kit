import React from "react"
import { Link, router, usePage } from "@inertiajs/react"
import { User, LogOut, ChevronDown } from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { PageProps } from "@/types"

export const UserMenu: React.FC = () => {
  const { auth } = usePage<PageProps>().props
  const user = auth?.user

  const handleLogout = () => {
    router.post("/logout")
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="flex items-center gap-2 pl-2 pr-3 py-1.5 h-10 hover:bg-slate-100 rounded-lg">
          <div className="h-8 w-8 rounded-full bg-indigo-100 border border-indigo-200 flex items-center justify-center text-indigo-700 font-semibold text-xs">
            {user?.name ? user.name.charAt(0).toUpperCase() : <User className="h-4 w-4" />}
          </div>
          <div className="hidden md:flex flex-col text-left">
            <span className="text-sm font-medium text-slate-800 leading-tight">{user?.name || "User"}</span>
            <span className="text-xs text-slate-500 leading-tight truncate max-w-[140px]">{user?.email}</span>
          </div>
          <ChevronDown className="h-4 w-4 text-slate-400" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel>
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">{user?.name}</p>
            <p className="text-xs leading-none text-muted-foreground">{user?.email}</p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <Link href="/dashboard" className="cursor-pointer">
            Dashboard
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <Link href="/products" className="cursor-pointer">
            Products
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <Link href="/users" className="cursor-pointer">
            Users
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onClick={handleLogout}
          className="text-red-600 focus:text-red-700 focus:bg-red-50 cursor-pointer"
        >
          <LogOut className="mr-2 h-4 w-4" />
          <span>Log out</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}