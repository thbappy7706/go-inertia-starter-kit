import { useState } from "react"
import { Head, router } from "@inertiajs/react"
import { AppLayout } from "@/layouts/AppLayout"
import { PageHeader } from "@/components/PageHeader"
import { Card, CardContent } from "@/components/ui/card"
import { ProductForm } from "./components/ProductForm"
import { ProductFormValues } from "@/lib/validations/product"
import { PageProps, Product } from "@/types"

interface ProductEditProps extends PageProps {
  product: Product
}

export default function ProductEdit({ product }: ProductEditProps) {
  const [loading, setLoading] = useState(false)

  const handleSubmit = (data: ProductFormValues) => {
    setLoading(true)
    router.put(`/products/${product.id}`, data, {
      onFinish: () => setLoading(false),
    })
  }

  return (
    <AppLayout title="Edit Product">
      <Head title={`Edit Product: ${product.name}`} />

      <div className="max-w-3xl mx-auto space-y-6">
        <PageHeader
          title={`Edit "${product.name}"`}
          description="Update details, pricing, and visibility for this product."
        />

        <Card className="border-slate-200/80 shadow-sm">
          <CardContent className="pt-6">
            <ProductForm
              product={product}
              onSubmit={handleSubmit}
              loading={loading}
            />
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  )
}