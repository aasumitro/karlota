import { createFileRoute } from '@tanstack/react-router'
import Terms from "@/features/terms";

export const Route = createFileRoute('/terms')({
  component: Terms,
})

