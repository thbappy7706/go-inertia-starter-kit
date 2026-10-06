import React, { useState } from "react"
import { Menu } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import { Sidebar } from "./Sidebar"
import { UserMenu } from "./UserMenu"

interface HeaderProps {
  title?: string
}

export const Header: React.FC<HeaderProps> = ({ title }) => {
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <header className="h-16 bg-white border-b border-slate-200 px-4 sm:px-8 flex items-center justify-between sticky top-0 z-30 shadow-xs">
      {/* Mobile drawer trigger + title */}
      <div className="flex items-center gap-3">
        <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon" className="lg:hidden text-slate-600">
              <Menu className="h-5 w-5" />
              <span className="sr-only">Open navigation menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="p-0 w-64 bg-slate-900 border-r border-slate-800">
            <Sidebar onItemClick={() => setMobileOpen(false)} />
          </SheetContent>
        </Sheet>

        {title && (
          <h2 className="text-lg font-semibold text-slate-800 hidden sm:block">
            {title}
          </h2>
        )}
      </div>

      {/* Right tools and User menu */}
      <div className="flex items-center gap-3">
        <UserMenu />
      </div>
    </header>
  )
}