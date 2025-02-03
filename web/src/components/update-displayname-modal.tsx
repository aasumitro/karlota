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
import {useUpdateDisplayName} from "@/hooks/use-account";
import {Profile} from "@/types/user";
import {useAuthStore} from "@/states/auth-store";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {Button} from "@/components/ui/button";
import {useGlobalActionStore} from "@/states/global-action-store";

export const UpdateDisplayNameModalState = "update_display_name_modal_state"

const formSchema = z.object({
  display_name: z
    .string()
    .min(1, { message: 'Please enter your display name.' }),
})

export function UpdateDisplayNameModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const { mutate: update, isPending} = useUpdateDisplayName();
  const {states, setState} = useGlobalActionStore();

  useEffect(() => setDialogOpen(states[UpdateDisplayNameModalState]), [states[UpdateDisplayNameModalState]]);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      display_name: useAuthStore.getState().auth.user?.display_name,
    },
  })

  const onSubmit = (data: z.infer<typeof formSchema>) => {
    update(data, {
      onSuccess: async (response) => {
        const resp = response?.data
        if (!resp) {
          return;
        }
        useAuthStore.getState().auth.setUser(resp as Profile)
        onClose();
      },
    });
  };

  const onClose = () => {
    setState(UpdateDisplayNameModalState, false);
    setDialogOpen(false);
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Update display name</AlertDialogTitle>
          <AlertDialogDescription />
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className='grid gap-2 mb-4'>
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