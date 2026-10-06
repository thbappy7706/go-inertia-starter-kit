import { Head, Link, usePage } from "@inertiajs/react"
import { AppLayout } from "@/layouts/AppLayout"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { StatusBadge } from "@/components/StatusBadge"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { PageProps, Product } from "@/types"
import { Users, Package, CheckCircle2, ArrowRight, Plus } from "lucide-react"
import { formatCurrency, formatDate } from "@/lib/utils"

interface DashboardStats {
  total_users: number
  total_products: number
  active_products: number
  recent_products: Product[]
}

interface DashboardProps extends PageProps {
  stats: DashboardStats
}

export default function Dashboard({ stats }: DashboardProps) {
  const { auth } = usePage<PageProps>().props
  const userName = auth?.user?.name || "Admin"

  return (
    <AppLayout title="Dashboard">
      <Head title="Dashboard" />

      <div className="space-y-8">
        {/* Welcome banner */}
        <div className="rounded-2xl bg-gradient-to-r from-indigo-700 via-indigo-600 to-indigo-800 p-6 sm:p-8 text-white shadow-lg relative overflow-hidden">
          <div className="relative z-10 max-w-2xl">
            <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight">
              Welcome back, {userName}!
            </h1>
            <p className="mt-2 text-indigo-100 text-sm sm:text-base leading-relaxed">
              Your Go + Inertia v3 application is running smoothly with PostgreSQL 17.
              Here is an overview of your platform metrics today.
            </p>
            <div className="mt-5 flex flex-wrap gap-3">
              <Button
                asChild
                className="bg-white text-indigo-700 hover:bg-indigo-50 font-semibold shadow-sm"
              >
                <Link href="/products/create">
                  <Plus className="mr-1.5 h-4 w-4" />
                  New Product
                </Link>
              </Button>
              <Button
                asChild
                variant="outline"
                className="border-indigo-400/40 text-white hover:bg-white/10 hover:text-white"
              >
                <Link href="/products">View All Products</Link>
              </Button>
            </div>
          </div>
          {/* Subtle decoration background circle */}
          <div className="absolute -right-12 -bottom-12 w-64 h-64 rounded-full bg-white/5 pointer-events-none" />
        </div>

        {/* Statistic Cards */}
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {/* Total Users */}
          <Card className="hover:shadow-md transition-shadow border-slate-200/80">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-500">
                Total Users
              </CardTitle>
              <div className="h-9 w-9 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
                <Users className="h-5 w-5" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold tracking-tight text-slate-900">
                {stats.total_users}
              </div>
              <p className="text-xs text-slate-500 mt-1 flex items-center gap-1">
                <span className="text-emerald-600 font-medium">Registered</span> accounts in database
              </p>
            </CardContent>
          </Card>

          {/* Total Products */}
          <Card className="hover:shadow-md transition-shadow border-slate-200/80">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-500">
                Total Products
              </CardTitle>
              <div className="h-9 w-9 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center">
                <Package className="h-5 w-5" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold tracking-tight text-slate-900">
                {stats.total_products}
              </div>
              <p className="text-xs text-slate-500 mt-1 flex items-center gap-1">
                Items currently cataloged
              </p>
            </CardContent>
          </Card>

          {/* Active Products */}
          <Card className="hover:shadow-md transition-shadow border-slate-200/80">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
              <CardTitle className="text-sm font-medium text-slate-500">
                Active Products
              </CardTitle>
              <div className="h-9 w-9 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center">
                <CheckCircle2 className="h-5 w-5" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold tracking-tight text-slate-900">
                {stats.active_products}
              </div>
              <p className="text-xs text-emerald-600 font-medium mt-1">
                Available for ordering
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Recent Products Table */}
        <Card className="border-slate-200/80 shadow-xs overflow-hidden">
          <CardHeader className="flex flex-row items-center justify-between pb-4">
            <div>
              <CardTitle className="text-lg font-bold text-slate-900">
                Recent Products
              </CardTitle>
              <CardDescription>
                Latest items added to the platform
              </CardDescription>
            </div>
            <Button asChild variant="ghost" size="sm" className="text-indigo-600 hover:text-indigo-700">
              <Link href="/products" className="inline-flex items-center gap-1">
                View all <ArrowRight className="h-4 w-4" />
              </Link>
            </Button>
          </CardHeader>
          <CardContent className="p-0">
            {stats.recent_products && stats.recent_products.length > 0 ? (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader className="bg-slate-50/75">
                    <TableRow>
                      <TableHead>Product Name</TableHead>
                      <TableHead>Price</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead className="text-right">Action</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {stats.recent_products.map((product) => (
                      <TableRow key={product.id} className="hover:bg-slate-50/50">
                        <TableCell className="font-semibold text-slate-900">
                          {product.name}
                          <span className="block text-xs font-normal text-slate-400 font-mono">
                            /{product.slug}
                          </span>
                        </TableCell>
                        <TableCell className="font-medium text-slate-700">
                          {formatCurrency(product.price)}
                        </TableCell>
                        <TableCell>
                          <StatusBadge status={product.status} />
                        </TableCell>
                        <TableCell className="text-slate-500 text-sm">
                          {formatDate(product.created_at)}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button asChild variant="outline" size="sm" className="h-8 text-xs">
                            <Link href={`/products/${product.id}/edit`}>Edit</Link>
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            ) : (
              <div className="p-8 text-center text-sm text-slate-500">
                No products found. Click &quot;New Product&quot; to create one.
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  )
}