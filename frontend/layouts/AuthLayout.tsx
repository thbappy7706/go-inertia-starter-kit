import { useEffect, ReactNode, FC } from "react"
import { usePage } from "@inertiajs/react"
import { Toaster, toast } from "sonner"
import { Layers } from "lucide-react"
import { PageProps } from "@/types"

interface AuthLayoutProps {
  children: ReactNode
  title: string
  subtitle: string
}

export const AuthLayout: FC<AuthLayoutProps> = ({
  children,
  title,
  subtitle,
}) => {
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
    <div className="min-h-screen flex flex-col justify-center py-12 sm:px-6 lg:px-8 bg-gradient-to-br from-slate-50 via-slate-100 to-indigo-50/40">
      <Toaster position="top-right" richColors closeButton />

      <div className="sm:mx-auto sm:w-full sm:max-w-md text-center">
        <div className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-indigo-600 text-white shadow-md shadow-indigo-500/20 mb-4">
          <Layers className="h-6 w-6" />
        </div>
        <h2 className="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
          {title}
        </h2>
        <p className="mt-2 text-sm text-slate-600">{subtitle}</p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md px-4 sm:px-0">
        <div className="bg-white py-8 px-6 shadow-xl shadow-slate-200/50 rounded-2xl border border-slate-100 sm:px-10">
          {children}
        </div>
      </div>
    </div>
  )
}