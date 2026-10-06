import { useState } from "react"
import { Head, Link, router, usePage } from "@inertiajs/react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { AuthLayout } from "@/layouts/AuthLayout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { registerSchema, RegisterFormValues } from "@/lib/validations/auth"
import { PageProps } from "@/types"
import { Loader2 } from "lucide-react"

export default function Register() {
  const { errors: serverErrors } = usePage<PageProps>().props
  const [loading, setLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      password_confirmation: "",
    },
  })

  const onSubmit = (data: RegisterFormValues) => {
    setLoading(true)
    router.post("/register", data, {
      onFinish: () => setLoading(false),
    })
  }

  return (
    <AuthLayout
      title="Create your account"
      subtitle="Join today and manage products seamlessly"
    >
      <Head title="Register" />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <Label htmlFor="name">Full Name</Label>
          <div className="mt-1">
            <Input
              id="name"
              type="text"
              autoComplete="name"
              placeholder="John Doe"
              {...register("name")}
            />
          </div>
          {(errors.name || serverErrors.name) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.name?.message || serverErrors.name}
            </p>
          )}
        </div>

        <div>
          <Label htmlFor="email">Email address</Label>
          <div className="mt-1">
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="john@example.com"
              {...register("email")}
            />
          </div>
          {(errors.email || serverErrors.email) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.email?.message || serverErrors.email}
            </p>
          )}
        </div>

        <div>
          <Label htmlFor="password">Password</Label>
          <div className="mt-1">
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              placeholder="At least 8 characters"
              {...register("password")}
            />
          </div>
          {(errors.password || serverErrors.password) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.password?.message || serverErrors.password}
            </p>
          )}
        </div>

        <div>
          <Label htmlFor="password_confirmation">Confirm Password</Label>
          <div className="mt-1">
            <Input
              id="password_confirmation"
              type="password"
              autoComplete="new-password"
              placeholder="Confirm password"
              {...register("password_confirmation")}
            />
          </div>
          {(errors.password_confirmation || serverErrors.password_confirmation) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.password_confirmation?.message || serverErrors.password_confirmation}
            </p>
          )}
        </div>

        <Button
          type="submit"
          className="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-2.5 font-medium mt-2"
          disabled={loading}
        >
          {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Create Account
        </Button>
      </form>

      <div className="mt-6 text-center text-sm text-slate-500">
        Already have an account?{" "}
        <Link
          href="/login"
          className="font-semibold text-indigo-600 hover:text-indigo-500"
        >
          Sign in
        </Link>
      </div>
    </AuthLayout>
  )
}