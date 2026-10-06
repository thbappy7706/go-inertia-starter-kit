import { useState } from "react"
import { Head, Link, router, usePage } from "@inertiajs/react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { AuthLayout } from "@/layouts/AuthLayout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { loginSchema, LoginFormValues } from "@/lib/validations/auth"
import { PageProps } from "@/types"
import { Loader2 } from "lucide-react"

export default function Login() {
  const { errors: serverErrors } = usePage<PageProps>().props
  const [loading, setLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
      remember: false,
    },
  })

  const onSubmit = (data: LoginFormValues) => {
    setLoading(true)
    router.post("/login", data, {
      onFinish: () => setLoading(false),
    })
  }

  return (
    <AuthLayout
      title="Welcome back"
      subtitle="Sign in to your account to access the dashboard"
    >
      <Head title="Login" />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div>
          <Label htmlFor="email">Email address</Label>
          <div className="mt-1">
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="admin@example.com"
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
          <div className="flex items-center justify-between">
            <Label htmlFor="password">Password</Label>
          </div>
          <div className="mt-1">
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••"
              {...register("password")}
            />
          </div>
          {(errors.password || serverErrors.password) && (
            <p className="mt-1 text-xs text-red-600">
              {errors.password?.message || serverErrors.password}
            </p>
          )}
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <input
              id="remember"
              type="checkbox"
              className="h-4 w-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
              {...register("remember")}
            />
            <Label htmlFor="remember" className="ml-2 text-sm text-slate-600 font-normal cursor-pointer">
              Remember me
            </Label>
          </div>
        </div>

        <Button
          type="submit"
          className="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-2.5 font-medium"
          disabled={loading}
        >
          {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Sign in
        </Button>
      </form>

      <div className="mt-6 text-center text-sm text-slate-500">
        Don&apos;t have an account?{" "}
        <Link
          href="/register"
          className="font-semibold text-indigo-600 hover:text-indigo-500"
        >
          Create account
        </Link>
      </div>

      <div className="mt-6 pt-4 border-t border-slate-100 text-xs text-slate-400 text-center">
        Development demo credentials: <br />
        <span className="font-mono text-slate-600 font-medium">admin@example.com</span> /{" "}
        <span className="font-mono text-slate-600 font-medium">password</span>
      </div>
    </AuthLayout>
  )
}