import React, { useState } from "react"
import { Head, Link, router } from "@inertiajs/react"
import { AppLayout } from "@/layouts/AppLayout"
import { PageHeader } from "@/components/PageHeader"
import { DataTable } from "@/components/DataTable"
import { Pagination } from "@/components/Pagination"
import { StatusBadge } from "@/components/StatusBadge"
import { EmptyState } from "@/components/EmptyState"
import { ConfirmDialog } from "@/components/ConfirmDialog"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { PageProps, PaginationMeta, Product } from "@/types"
import { Plus, Search, MoreHorizontal, Edit, Trash2, Package, X } from "lucide-react"
import { formatCurrency, formatDate } from "@/lib/utils"

interface ProductsIndexProps extends PageProps {
  products: PaginationMeta<Product>
  filters: {
    search?: string
    status?: string
  }
}

export default function ProductsIndex({ products, filters }: ProductsIndexProps) {
  const [searchTerm, setSearchTerm] = useState(filters.search || "")
  const [statusFilter, setStatusFilter] = useState(filters.status || "all")
  const [productToDelete, setProductToDelete] = useState<Product | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const applyFilters = (newSearch?: string, newStatus?: string) => {
    const s = newSearch !== undefined ? newSearch : searchTerm
    const st = newStatus !== undefined ? newStatus : statusFilter

    const query: Record<string, string | number> = { page: 1 }
    if (s.trim()) query.search = s.trim()
    if (st && st !== "all") query.status = st

    router.get("/products", query, { preserveState: true })
  }

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    applyFilters(searchTerm, statusFilter)
  }

  const handleStatusChange = (val: string) => {
    setStatusFilter(val)
    applyFilters(searchTerm, val)
  }

  const handleClear = () => {
    setSearchTerm("")
    setStatusFilter("all")
    router.get("/products", {}, { preserveState: true })
  }

  const handleDeleteConfirm = () => {
    if (!productToDelete) return
    router.delete(`/products/${productToDelete.id}`, {
      onSuccess: () => {
        setDeleteDialogOpen(false)
        setProductToDelete(null)
      },
    })
  }

  return (
    <AppLayout title="Products">
      <Head title="Products" />

      <div className="space-y-6">
        <PageHeader
          title="Products"
          description="Manage your item catalog, pricing, and availability."
        >
          <Button asChild className="bg-indigo-600 hover:bg-indigo-700 text-white shadow-sm">
            <Link href="/products/create">
              <Plus className="mr-1.5 h-4 w-4" />
              Add Product
            </Link>
          </Button>
        </PageHeader>

        {/* Filter Toolbar */}
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4">
          <form onSubmit={handleSearchSubmit} className="flex items-center gap-3 flex-1 max-w-md">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
              <Input
                type="text"
                placeholder="Search products..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-9 h-9"
              />
              {searchTerm && (
                <button
                  type="button"
                  onClick={() => {
                    setSearchTerm("")
                    applyFilters("", statusFilter)
                  }}
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

          <div className="flex items-center gap-3">
            <div className="w-36">
              <Select value={statusFilter} onValueChange={handleStatusChange}>
                <SelectTrigger className="h-9">
                  <SelectValue placeholder="All Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {(filters.search || (filters.status && filters.status !== "all")) && (
              <Button variant="ghost" size="sm" onClick={handleClear} className="text-slate-500 h-9">
                Clear
              </Button>
            )}
          </div>
        </div>

        {/* Products Table */}
        <DataTable>
          {products.data && products.data.length > 0 ? (
            <div>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader className="bg-slate-50">
                    <TableRow>
                      <TableHead>Product</TableHead>
                      <TableHead>Price</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead className="text-right">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {products.data.map((product) => (
                      <TableRow key={product.id} className="hover:bg-slate-50/50">
                        <TableCell>
                          <div className="font-semibold text-slate-900">
                            {product.name}
                          </div>
                          <div className="text-xs font-mono text-slate-400">
                            {product.slug}
                          </div>
                        </TableCell>
                        <TableCell className="font-semibold text-slate-700">
                          {formatCurrency(product.price)}
                        </TableCell>
                        <TableCell>
                          <StatusBadge status={product.status} />
                        </TableCell>
                        <TableCell className="text-slate-500 text-sm">
                          {formatDate(product.created_at)}
                        </TableCell>
                        <TableCell className="text-right">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon" className="h-8 w-8">
                                <MoreHorizontal className="h-4 w-4" />
                                <span className="sr-only">Actions</span>
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end" className="w-40">
                              <DropdownMenuItem asChild>
                                <Link
                                  href={`/products/${product.id}/edit`}
                                  className="cursor-pointer flex items-center"
                                >
                                  <Edit className="mr-2 h-4 w-4 text-slate-500" />
                                  Edit
                                </Link>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => {
                                  setProductToDelete(product)
                                  setDeleteDialogOpen(true)
                                }}
                                className="cursor-pointer text-red-600 focus:text-red-700 focus:bg-red-50 flex items-center"
                              >
                                <Trash2 className="mr-2 h-4 w-4" />
                                Delete
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              <div className="border-t border-slate-100">
                <Pagination
                  currentPage={products.current_page}
                  lastPage={products.last_page}
                  total={products.total}
                  from={products.from}
                  to={products.to}
                />
              </div>
            </div>
          ) : (
            <div className="p-8">
              <EmptyState
                icon={Package}
                title="No products found"
                description={
                  filters.search || filters.status
                    ? "Try adjusting your search criteria or status filter."
                    : "Get started by creating your first product."
                }
                action={
                  <Button asChild className="bg-indigo-600 hover:bg-indigo-700 text-white">
                    <Link href="/products/create">
                      <Plus className="mr-1.5 h-4 w-4" />
                      Add Product
                    </Link>
                  </Button>
                }
              />
            </div>
          )}
        </DataTable>
      </div>

      {/* Delete Confirmation Alert Dialog */}
      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Delete Product?"
        description={`Are you sure you want to delete "${productToDelete?.name}"? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDeleteConfirm}
        destructive
      />
    </AppLayout>
  )
}