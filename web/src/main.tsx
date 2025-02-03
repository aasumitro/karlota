import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { AxiosError } from 'axios'
import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { useAuthStore } from '@/states/auth-store'
import {handleServerError, HTTP_STATUS_CODE} from '@/lib/api'
import { toast } from '@/hooks/use-toast'
import './index.css'
// Generated Routes
import { routeTree } from './routeTree.gen'

const basePath = process.env.NODE_ENV === 'production' ? '/ui' : '/'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        if (import.meta.env.DEV) console.log({ failureCount, error })
        if (failureCount >= 0 && import.meta.env.DEV) return false
        if (failureCount > 3 && import.meta.env.PROD) return false
        return !(
          error instanceof AxiosError &&
          [
            HTTP_STATUS_CODE.INTERNAL_SERVER_ERROR,
            HTTP_STATUS_CODE.FORBIDDEN
          ].includes(error.response?.status ?? 0)
        )
      },
      refetchOnWindowFocus: import.meta.env.PROD,
      staleTime: 10 * 1000, // 10s
    },
    mutations: {
      onError: (error) => {
        handleServerError(error)
        if (error instanceof AxiosError) {
          if (error.response?.status === HTTP_STATUS_CODE.UNAUTHORIZED) {
            toast({
              variant: 'destructive',
              title: 'Session expired!',
            })
            useAuthStore.getState().auth.reset()
            const redirect = `${router.history.location.href}`
            router.navigate({ to: '/sign-in', search: { redirect } })
          }
          if (error.response?.status === HTTP_STATUS_CODE.NOT_MODIFIED) {
            toast({
              variant: 'destructive',
              title: 'Content not modified!',
            })
          }
        }
      },
    },
  },
  queryCache: new QueryCache({
    onError: (error) => {
      if (error instanceof AxiosError) {
        if (error.response?.status === HTTP_STATUS_CODE.UNAUTHORIZED) {
          toast({
            variant: 'destructive',
            title: 'Session expired!',
          })
          useAuthStore.getState().auth.reset()
          const redirect = `${router.history.location.href}`
          router.navigate({ to: '/sign-in', search: { redirect } })
        }
        if (error.response?.status === HTTP_STATUS_CODE.INTERNAL_SERVER_ERROR) {
          toast({
            variant: 'destructive',
            title: 'Internal Server Error!',
          })
          router.navigate({ to: '/500' })
        }
      }
    },
  }),
})

// Create a new router instance
const router = createRouter({
  routeTree,
  context: { queryClient },
  defaultPreload: 'intent',
  defaultPreloadStaleTime: 0,
  basepath: basePath
})

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  // noinspection JSUnusedGlobalSymbols
  interface Register {
    router: typeof router
  }
}

// Render the app
const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </StrictMode>
  )
}
