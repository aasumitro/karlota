import { createFileRoute } from '@tanstack/react-router'
import Privacy from "@/features/privacy";

export const Route = createFileRoute('/privacy')({
  component: Privacy,
})
