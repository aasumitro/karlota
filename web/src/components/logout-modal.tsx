import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "./ui/alert-dialog";
import {useState, FormEvent, useEffect} from "react";
import { Loader} from "lucide-react";
import { Button } from "./ui/button";
import {useAuthStore} from "@/states/auth-store";
import {useNavigate} from "@tanstack/react-router";
import {useGlobalActionStore} from "@/states/global-action-store";

export const LogoutModalState = "logout_modal_state"

export function LogoutModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const [isSubmitted, setSubmitted] = useState(false);
  const navigate = useNavigate();
  const { states, setState } = useGlobalActionStore();

  useEffect(() => {
    if (states[LogoutModalState]){
      setDialogOpen(states[LogoutModalState])
    }
  }, [states]);

  const onSubmit = (event: FormEvent) => {
    event.preventDefault();
    setSubmitted(true);
    setTimeout(async () => {
      onClose();
      useAuthStore.getState().auth.reset();
      await navigate({to: "/sign-in", replace: true})
    }, 1000);
  };

  const onClose = () => {
    setState(LogoutModalState, false)
    setDialogOpen(false)
    setSubmitted(false);
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Confirm Logout</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to log out? You will need to re-enter your credentials to access your account
            again. Any unsaved changes will not be saved and progress might be lost.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <form onSubmit={onSubmit}>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={onClose}>
              Cancel
            </AlertDialogCancel>
            <Button
              className="bg-red-500 hover:bg-red-600 text-white"
              type="submit"
              disabled={isSubmitted}
            >
              <Loader
                className={`${
                  isSubmitted ? "block animate-spin mr-2" : "hidden"
                }`}
              />
              Confirm
            </Button>
          </AlertDialogFooter>
        </form>
      </AlertDialogContent>
    </AlertDialog>
  );
}