import { useEffect, ReactNode, FC } from "react"
import { usePage } from "@inertiajs/react"
import { Toaster, toast } from "sonner"
import { Sidebar } from "./components/Sidebar"
import { Header } from "./components/Header"
import { PageProps } from "@/types"

interface AppLayoutProps {
  title?: string
  children: ReactNode
}

export const AppLayout: FC<AppLayoutProps> = ({ title, children }) => {
  const { flash } = usePage<PageProps>().props

  useEffect(() => {
    if (flash?.success) {
      toast.success(flash.success)
    }
    if (flash?.error) {
      toast.error(flash.error)
    }
  }, [flash])

  return (
    <div className="min-h-screen flex bg-slate-50 font-sans text-slate-900">
      <Toaster position="top-right" richColors closeButton />

      {/* Desktop Sidebar */}
      <div className="hidden lg:block shrink-0">
        <div className="sticky top-0 h-screen">
          <Sidebar />
        </div>
      </div>

      {/* Main body area */}
      <div className="flex-1 flex flex-col min-w-0">
        <Header title={title} />
        <main className="flex-1 p-4 sm:p-6 lg:p-8 max-w-7xl w-full mx-auto">
          {children}
        </main>
      </div>
    </div>
  )
}