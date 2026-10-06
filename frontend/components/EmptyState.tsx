import React from "react"
import { LucideIcon, FolderSearch } from "lucide-react"

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description?: string
  action?: React.ReactNode
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  icon: Icon = FolderSearch,
  title,
  description,
  action,
}) => {
  return (
    <div className="flex flex-col items-center justify-center p-8 sm:p-12 text-center rounded-xl border border-dashed border-slate-200 bg-slate-50/50">
      <div className="h-12 w-12 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 mb-4 shadow-sm">
        <Icon className="h-6 w-6" />
      </div>
      <h3 className="text-base font-semibold text-slate-800">{title}</h3>
      {description && (
        <p className="mt-1 text-sm text-slate-500 max-w-sm mb-6">{description}</p>
      )}
      {action && <div>{action}</div>}
    </div>
  )
}