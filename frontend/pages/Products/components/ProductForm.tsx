import React, { useState } from "react"
import { Link, usePage } from "@inertiajs/react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { productSchema, ProductFormValues } from "@/lib/validations/product"
import { PageProps, Product } from "@/types"
import { Loader2 } from "lucide-react"

interface ProductFormProps {
  product?: Product
  onSubmit: (data: ProductFormValues) => void
  loading?: boolean
}

export const ProductForm: React.FC<ProductFormProps> = ({
  product,
  onSubmit,
  loading = false,
}) => {
  const { errors: serverErrors } = usePage<PageProps>().props
  const [autoSlug, setAutoSlug] = useState(!product)

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      name: product?.name || "",
      slug: product?.slug || "",
      description: product?.description || "",
      price: product ? Number(product.price) : 0,
      status: product ? product.status : true,
    },
  })

  const generateSlug = (text: string) => {
    return text
      .toLowerCase()
      .replace(/[^\w\s-]/g, "")
      .replace(/[\s_-]+/g, "-")
      .replace(/^-+|-+$/g, "")
  }

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    if (autoSlug) {
      setValue("slug", generateSlug(val), { shouldValidate: true })
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        {/* Name */}
        <div className="sm:col-span-2">
          <Label htmlFor="name">Product Name *</Label>
          <div className="mt-1.5">
            <Input
              id="name"
              placeholder="e.g. Enterprise Cloud Hosting"
              {...register("name", {
                onChange: handleNameChange,
              })}
            />
          </div>
          {(errors.name || serverErrors.name) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.name?.message || serverErrors.name}
            </p>
          )}
        </div>

        {/* Slug */}
        <div className="sm:col-span-2">
          <div className="flex items-center justify-between">
            <Label htmlFor="slug">Slug (URL identifier) *</Label>
            {!product && (
              <button
                type="button"
                onClick={() => setAutoSlug(!autoSlug)}
                className="text-xs text-indigo-600 hover:text-indigo-700"
              >
                {autoSlug ? "Manual slug" : "Auto generate"}
              </button>
            )}
          </div>
          <div className="mt-1.5 flex rounded-lg shadow-xs">
            <span className="inline-flex items-center px-3 rounded-l-lg border border-r-0 border-input bg-slate-50 text-slate-500 text-xs font-mono">
              /products/
            </span>
            <Input
              id="slug"
              className="rounded-l-none font-mono text-xs"
              placeholder="enterprise-cloud-hosting"
              {...register("slug")}
              onChange={() => setAutoSlug(false)}
            />
          </div>
          {(errors.slug || serverErrors.slug) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.slug?.message || serverErrors.slug}
            </p>
          )}
        </div>

        {/* Price */}
        <div>
          <Label htmlFor="price">Price (USD) *</Label>
          <div className="mt-1.5 relative rounded-lg shadow-xs">
            <span className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-500 text-sm">
              $
            </span>
            <Input
              id="price"
              type="number"
              step="0.01"
              min="0"
              className="pl-7"
              placeholder="29.99"
              {...register("price")}
            />
          </div>
          {(errors.price || serverErrors.price) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.price?.message || serverErrors.price}
            </p>
          )}
        </div>

        {/* Status */}
        <div>
          <Label htmlFor="status">Product Status</Label>
          <div className="mt-1.5 flex items-center h-9">
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                id="status"
                type="checkbox"
                className="sr-only peer"
                {...register("status")}
              />
              <div className="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-emerald-600" />
              <span className="ml-3 text-sm font-medium text-slate-700">
                Active in catalog
              </span>
            </label>
          </div>
        </div>

        {/* Description */}
        <div className="sm:col-span-2">
          <Label htmlFor="description">Description (Optional)</Label>
          <div className="mt-1.5">
            <Textarea
              id="description"
              rows={4}
              placeholder="Enter detailed description of features, specifications, etc..."
              {...register("description")}
            />
          </div>
          {errors.description && (
            <p className="mt-1 text-xs text-red-600">
              {errors.description?.message}
            </p>
          )}
        </div>
      </div>

      <div className="flex items-center justify-end gap-3 pt-6 border-t border-slate-100">
        <Button asChild variant="outline">
          <Link href="/products">Cancel</Link>
        </Button>
        <Button
          type="submit"
          disabled={loading}
          className="bg-indigo-600 hover:bg-indigo-700 text-white min-w-[120px]"
        >
          {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          {product ? "Update Product" : "Create Product"}
        </Button>
      </div>
    </form>
  )
}