import {HTMLAttributes, useState} from 'react'
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
import {PasswordInput} from "@/components/password-input";
import {useResetPassword} from "@/features/auth/reset-password/hooks/use-reset-password";
import {useNavigate} from "@tanstack/react-router";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from "@/components/ui/alert-dialog";

type ResetFormProps = HTMLAttributes<HTMLDivElement>

const formSchema = z.object({
  new_password: z
    .string()
    .min(1, {
      message: 'Please enter your password',
    })
    .min(7, {
      message: 'Password must be at least 7 characters long',
    }),
  confirm_password: z.string(),
  exchange_token: z.string(),
}).refine((data) => data.new_password === data.confirm_password, {
  message: "Passwords don't match.",
  path: ['confirmPassword'],
})

export function ResetPasswordForm({ className, ...props }: ResetFormProps) {
  const queryParams = new URLSearchParams(location.search);
  const exchangeToken = queryParams.get('exchange_token');
  const navigate = useNavigate();
  let [openModal, setOpenModal] = useState(false);

  if (!exchangeToken) {
    return ;
  }

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      new_password: '',
      confirm_password: '',
      exchange_token: exchangeToken as string
    },
  })

  const {mutate: resetPassword, isPending} = useResetPassword();

  function onSubmit(data: z.infer<typeof formSchema>) {
    resetPassword(data, {
      onSuccess: async (response) => {
        const resp = response?.data
        if (!resp) {
          return;
        }
        setOpenModal(true);
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
              name='new_password'
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
              {isPending ? "Please wait..." : "Continue"}
            </Button>
          </div>
        </form>
      </Form>

      <AlertDialog open={openModal} >
        <AlertDialogContent onEscapeKeyDown={(e) => e.preventDefault()}>
          <AlertDialogHeader>
            <AlertDialogTitle>Action Request Success</AlertDialogTitle>
            <AlertDialogDescription>
              Your password has been successfully changed.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button
              className="text-white"
              onClick={async () => await navigate({to: '/sign-in', replace: true})}
            >
              Sign in
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}