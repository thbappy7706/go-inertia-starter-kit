import "./css/app.css"
import { createRoot } from "react-dom/client"
import { createInertiaApp } from "@inertiajs/react"

createInertiaApp({
  title: (title) => (title ? `${title} - Go Inertia Starter Kit` : "Go Inertia Starter Kit"),
  resolve: (name) => {
    const pages = import.meta.glob("./pages/**/*.tsx", { eager: true })
    const page = (pages as Record<string, any>)[`./pages/${name}.tsx`]
    if (!page) {
      throw new Error(`Page not found: ${name}`)
    }
    return page
  },
  setup({ el, App, props }) {
    createRoot(el).render(<App {...props} />)
  },
})