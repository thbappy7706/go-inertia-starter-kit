import React from "react"
import { router } from "@inertiajs/react"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"

interface PaginationProps {
  currentPage: number
  lastPage: number
  total: number
  from: number
  to: number
  onPageChange?: (page: number) => void
  preserveScroll?: boolean
}

export const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  lastPage,
  total,
  from,
  to,
  onPageChange,
  preserveScroll = true,
}) => {
  if (total === 0) return null

  const handlePageClick = (page: number) => {
    if (page < 1 || page > lastPage || page === currentPage) return
    if (onPageChange) {
      onPageChange(page)
    } else {
      const url = new URL(window.location.href)
      url.searchParams.set("page", page.toString())
      router.get(url.pathname + url.search, {}, { preserveState: true, preserveScroll })
    }
  }

  // Generate page numbers
  const pages: (number | string)[] = []
  if (lastPage <= 7) {
    for (let i = 1; i <= lastPage; i++) pages.push(i)
  } else {
    pages.push(1)
    if (currentPage > 3) pages.push("...")
    const start = Math.max(2, currentPage - 1)
    const end = Math.min(lastPage - 1, currentPage + 1)
    for (let i = start; i <= end; i++) pages.push(i)
    if (currentPage < lastPage - 2) pages.push("...")
    pages.push(lastPage)
  }

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 px-2 py-4">
      <div className="text-sm text-muted-foreground">
        Showing <span className="font-semibold text-foreground">{from}</span> to{" "}
        <span className="font-semibold text-foreground">{to}</span> of{" "}
        <span className="font-semibold text-foreground">{total}</span> results
      </div>

      <div className="flex items-center space-x-1.5">
        <Button
          variant="outline"
          size="sm"
          onClick={() => handlePageClick(currentPage - 1)}
          disabled={currentPage <= 1}
          className="h-8 w-8 p-0"
        >
          <ChevronLeft className="h-4 w-4" />
          <span className="sr-only">Previous page</span>
        </Button>

        {pages.map((p, idx) => {
          if (p === "...") {
            return (
              <span key={`dots-${idx}`} className="px-2 text-xs text-muted-foreground">
                ...
              </span>
            )
          }

          const pageNum = p as number
          const isActive = pageNum === currentPage
          return (
            <Button
              key={pageNum}
              variant={isActive ? "default" : "outline"}
              size="sm"
              onClick={() => handlePageClick(pageNum)}
              className={`h-8 w-8 p-0 text-xs ${
                isActive ? "bg-indigo-600 text-white hover:bg-indigo-700" : ""
              }`}
            >
              {pageNum}
            </Button>
          )
        })}

        <Button
          variant="outline"
          size="sm"
          onClick={() => handlePageClick(currentPage + 1)}
          disabled={currentPage >= lastPage}
          className="h-8 w-8 p-0"
        >
          <ChevronRight className="h-4 w-4" />
          <span className="sr-only">Next page</span>
        </Button>
      </div>
    </div>
  )
}