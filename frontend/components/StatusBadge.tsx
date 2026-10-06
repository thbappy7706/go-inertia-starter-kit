import React from "react"
import { Badge } from "@/components/ui/badge"

interface StatusBadgeProps {
  status: boolean
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({ status }) => {
  return (
    <Badge
      variant={status ? "success" : "inactive"}
      className="inline-flex items-center gap-1.5 font-medium px-2.5 py-0.5"
    >
      <span
        className={`h-1.5 w-1.5 rounded-full ${
          status ? "bg-emerald-500" : "bg-slate-400"
        }`}
      />
      {status ? "Active" : "Inactive"}
    </Badge>
  )
}