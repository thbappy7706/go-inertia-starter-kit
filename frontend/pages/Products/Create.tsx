import { useState } from "react"
import { Head, router } from "@inertiajs/react"
import { AppLayout } from "@/layouts/AppLayout"
import { PageHeader } from "@/components/PageHeader"
import { Card, CardContent } from "@/components/ui/card"
import { ProductForm } from "./components/ProductForm"
import { ProductFormValues } from "@/lib/validations/product"

export default function ProductCreate() {
  const [loading, setLoading] = useState(false)

  const handleSubmit = (data: ProductFormValues) => {
    setLoading(true)
    router.post("/products", data, {
      onFinish: () => setLoading(false),
    })
  }

  return (
    <AppLayout title="Create Product">
      <Head title="Create Product" />

      <div className="max-w-3xl mx-auto space-y-6">
        <PageHeader
          title="Create New Product"
          description="Add a new catalog item with pricing and status details."
        />

        <Card className="border-slate-200/80 shadow-sm">
          <CardContent className="pt-6">
            <ProductForm onSubmit={handleSubmit} loading={loading} />
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  )
}