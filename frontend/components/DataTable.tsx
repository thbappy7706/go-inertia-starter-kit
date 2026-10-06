import React from "react"

interface DataTableProps {
  children: React.ReactNode
}

export const DataTable: React.FC<DataTableProps> = ({ children }) => {
  return (
    <div className="rounded-xl border border-slate-200 bg-white shadow-sm overflow-hidden">
      {children}
    </div>
  )
}