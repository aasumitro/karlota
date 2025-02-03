import {
  AlertDialog, AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "./ui/alert-dialog";
import {useEffect, useState} from "react";
import {Loader} from "lucide-react";
import {z} from "zod";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {useUpdatePassword} from "@/hooks/use-account";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Button} from "@/components/ui/button";
import {PasswordInput} from "@/components/password-input";
import {useAuthStore} from "@/states/auth-store";
import {useNavigate} from "@tanstack/react-router";
import {useGlobalActionStore} from "@/states/global-action-store";

export const UpdatePasswordModalState = "update_password_modal_state"

const formSchema = z.object({
  old_password: z
    .string()
    .min(1, {
      message: 'Please enter your password',
    })
    .min(7, {
      message: 'Password must be at least 7 characters long',
    }),
  new_password: z
    .string()
    .min(1, {
      message: 'Please enter your password',
    })
    .min(7, {
      message: 'Password must be at least 7 characters long',
    }),
  new_confirm_password: z.string(),
}).refine((data) => data.new_password === data.new_confirm_password, {
  message: "Passwords don't match.",
  path: ['new_confirm_password'],
})

export function UpdatePasswordModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { mutate: update, isPending} = useUpdatePassword();
  const navigate = useNavigate();
  const {states, setState} = useGlobalActionStore();

  useEffect(() => setDialogOpen(states[UpdatePasswordModalState]), [states[UpdatePasswordModalState]]);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      old_password: '',
      new_password: '',
      new_confirm_password: '',
    },
  })

  const onSubmit = (data: z.infer<typeof formSchema>) => {
    update(data, {
      onSuccess: async (response) => {
        const resp = response?.data
        if (!resp) {
          return;
        }
        useAuthStore.getState().auth.reset();
        await navigate({to: "/sign-in", replace: true})
        onClose();
      },
    });
  };

  const onClose = () => {
    setState(UpdatePasswordModalState, false);
    setDialogOpen(false);
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Update password</AlertDialogTitle>
          <AlertDialogDescription />
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className='grid gap-2 mb-4'>
              <FormField
                control={form.control}
                name='old_password'
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
                name='new_confirm_password'
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
            </div>
            <AlertDialogFooter>
              <AlertDialogCancel onClick={onClose}>
                Cancel
              </AlertDialogCancel>
              <Button
                className="text-white"
                type="submit"
              >
                <Loader
                  className={`${
                    isPending ? "block animate-spin mr-2" : "hidden"
                  }`}
                />
                Confirm
              </Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}