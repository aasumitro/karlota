import { HTMLAttributes } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/password-input'
import {useAuthStore} from "@/states/auth-store";
import {useNavigate} from "@tanstack/react-router";
import {useRegister} from "@/features/auth/sign-up/hooks/use-register";
import {Profile} from "@/types/user";

type SignUpFormProps = HTMLAttributes<HTMLDivElement>

const formSchema = z
  .object({
    display_name: z
      .string()
      .min(1, { message: 'Please enter your display name.' }),
    email: z
      .string()
      .min(1, { message: 'Please enter your email' })
      .email({ message: 'Invalid email address' }),
    password: z
      .string()
      .min(1, {
        message: 'Please enter your password',
      })
      .min(7, {
        message: 'Password must be at least 7 characters long',
      }),
    confirm_password: z.string(),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords don't match.",
    path: ['confirmPassword'],
  })

export function SignUpForm({ className, ...props }: SignUpFormProps) {
  const {auth} = useAuthStore();
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      display_name: '',
      email: '',
      password: '',
      confirm_password: '',
    },
  })

  const { mutate: register, isPending} = useRegister();

  function onSubmit(data: z.infer<typeof formSchema>) {
    register(data, {
      onSuccess: async (response) => {
        const resp = response?.data
        if (!resp) {
          return;
        }
        auth.setAccessToken(resp.access_token.str as string)
        auth.setRefreshToken(resp.refresh_token.str as string)
        auth.setUser(resp.user as Profile)
        await navigate({to: '/chats', replace: true})
      },
    });
  }

  return (
    <div className={cn('grid gap-6', className)} {...props}>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <div className='grid gap-2'>
            <FormField
              control={form.control}
              name='display_name'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>Display name</FormLabel>
                  <FormControl>
                    <Input placeholder='@itsmebro' {...field} />
                  </FormControl>
                  <FormMessage />
                  <p className={"text-xs m-1 text-slate-500 dark:text-slate-400"}>
                    It will be visible to other users. Do not use the '@' symbol,
                    just type your name. For example: 'lorem'`
                  </p>
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='email'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder='name@example.com' {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <PasswordInput placeholder='********' {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='confirm_password'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>Confirm Password</FormLabel>
                  <FormControl>
                    <PasswordInput placeholder='********' {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button className='mt-2' disabled={isPending}>
              {isPending ? "Please wait..." : "Create Account"}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  )
}