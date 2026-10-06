import React, { useState } from "react"
import { Head, router } from "@inertiajs/react"
import { AppLayout } from "@/layouts/AppLayout"
import { PageHeader } from "@/components/PageHeader"
import { DataTable } from "@/components/DataTable"
import { Pagination } from "@/components/Pagination"
import { EmptyState } from "@/components/EmptyState"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { PageProps, PaginationMeta, User } from "@/types"
import { Search, Users, X } from "lucide-react"
import { formatDate } from "@/lib/utils"

interface UsersIndexProps extends PageProps {
  users: PaginationMeta<User>
  filters: {
    search?: string
  }
}

export default function UsersIndex({ users, filters }: UsersIndexProps) {
  const [searchTerm, setSearchTerm] = useState(filters.search || "")

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    router.get(
      "/users",
      { search: searchTerm, page: 1 },
      { preserveState: true }
    )
  }

  const handleClear = () => {
    setSearchTerm("")
    router.get("/users", {}, { preserveState: true })
  }

  return (
    <AppLayout title="Users">
      <Head title="Users" />

      <div className="space-y-6">
        <PageHeader
          title="Users"
          description="View all registered users and system administrators."
        />

        {/* Filter Toolbar */}
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-xs">
          <form onSubmit={handleSearch} className="flex items-center gap-3 max-w-md">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
              <Input
                type="text"
                placeholder="Search users by name or email..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-9 h-9"
              />
              {searchTerm && (
                <button
                  type="button"
                  onClick={handleClear}
                  className="absolute right-2.5 top-2.5 text-slate-400 hover:text-slate-600"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>
            <Button type="submit" size="sm" className="bg-indigo-600 hover:bg-indigo-700 h-9">
              Search
            </Button>
          </form>
        </div>

        {/* Users Table */}
        <DataTable>
          {users.data && users.data.length > 0 ? (
            <div>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader className="bg-slate-50">
                    <TableRow>
                      <TableHead>User</TableHead>
                      <TableHead>Email Address</TableHead>
                      <TableHead>Joined Date</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {users.data.map((user) => (
                      <TableRow key={user.id} className="hover:bg-slate-50/50">
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <div className="h-8 w-8 rounded-full bg-indigo-100 text-indigo-700 flex items-center justify-center font-bold text-xs uppercase">
                              {user.name.charAt(0)}
                            </div>
                            <span className="font-semibold text-slate-900">
                              {user.name}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-slate-600">
                          {user.email}
                        </TableCell>
                        <TableCell className="text-slate-500 text-sm">
                          {formatDate(user.created_at)}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              <div className="border-t border-slate-100">
                <Pagination
                  currentPage={users.current_page}
                  lastPage={users.last_page}
                  total={users.total}
                  from={users.from}
                  to={users.to}
                />
              </div>
            </div>
          ) : (
            <div className="p-8">
              <EmptyState
                icon={Users}
                title="No users found"
                description={
                  filters.search
                    ? "Try adjusting your search criteria to find registered accounts."
                    : "No users are registered in the system yet."
                }
                action={
                  filters.search ? (
                    <Button variant="outline" size="sm" onClick={handleClear}>
                      Clear filter
                    </Button>
                  ) : null
                }
              />
            </div>
          )}
        </DataTable>
      </div>
    </AppLayout>
  )
}