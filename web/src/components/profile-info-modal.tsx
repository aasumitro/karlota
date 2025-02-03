import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "./ui/alert-dialog";
import {useEffect, useState} from "react";
import {Phone, Video} from "lucide-react";
import {Button} from "@/components/ui/button";
import {Avatar, AvatarFallback} from "@/components/ui/avatar";
import {useAuthStore} from "@/states/auth-store";
import {useGlobalActionStore} from "@/states/global-action-store";
import {useSelectActionStore} from "@/states/select-action-store";
import {formatOnlineTime} from "@/lib/time";

export const ProfileInfoModalState = "profile_info_modal_state"

export function ProfileInfoModal() {
  const [isDialogOpen, setDialogOpen] = useState(false);
  const {auth} = useAuthStore();
  const {states, setState} = useGlobalActionStore();
  const {user, setUser} = useSelectActionStore();

  useEffect(() => setDialogOpen(states[ProfileInfoModalState]), [states[ProfileInfoModalState]]);

  const onClose = () => {
    setState(ProfileInfoModalState, false);
    setDialogOpen(false);
    setUser(null);
  }

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle className="flex justify-between items-center">
            Info
            <Button variant="ghost" onClick={onClose}>X</Button>
          </AlertDialogTitle>
          <AlertDialogDescription/>
        </AlertDialogHeader>
        <section className="flex flex-col items-center justify-center">
          <Avatar className='size-11 lg:size-14'>
            <AvatarFallback>
              {user?.display_name[0] ?? "U"}
            </AvatarFallback>
          </Avatar>
          <div className="text-center mt-4">
            <p className="font-semibold text-2xl mb-0.5">
              {user?.display_name}
            </p>
            <p className="font-medium text-md text-gray-500">
              {user?.email ? user.email :
                user?.is_online
                  ? "online"
                  : formatOnlineTime(user?.last_online)}
            </p>
          </div>
          <div className="flex flex-row items-center gap-4 mt-4">
            <Button
              variant="outline"
              className="flex flex-col h-16 w-30 text-[10px]"
              disabled={auth?.user?.id == user?.id}
            >
              <Phone className="w-4 h-4" />
              Audio
            </Button>
            <Button
              variant="outline"
              className="flex flex-col h-16 w-30 text-[10px]"
              disabled={auth?.user?.id == user?.id}
            >
              <Video className="w-4 h-4" />
              Video
            </Button>
          </div>
        </section>
      </AlertDialogContent>
    </AlertDialog>
  );
}