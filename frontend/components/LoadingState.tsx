import React from "react"
import { Skeleton } from "@/components/ui/skeleton"

export const LoadingState: React.FC = () => {
  return (
    <div className="space-y-4 p-4">
      <Skeleton className="h-8 w-1/3" />
      <div className="space-y-2">
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-12 w-full" />
      </div>
    </div>
  )
}