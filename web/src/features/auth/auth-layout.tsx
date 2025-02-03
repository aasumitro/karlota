import {IconMessages} from "@tabler/icons-react";
import {useAuthStore} from "@/states/auth-store";
import {isJWTExpired} from "@/lib/jwt";
import {toast} from "@/hooks/use-toast";
import {useNavigate} from "@tanstack/react-router";
import {ReactNode} from "react";

interface Props {
  children: ReactNode
}

export default function AuthLayout({ children }: Props) {
  const navigate = useNavigate()
  if (useAuthStore.getState().auth.accessToken &&
    !isJWTExpired(useAuthStore.getState().auth.accessToken)) {
    navigate({to: "/chats", replace: true}).then(() =>
      toast({ variant: 'default', title: "session valid" }))
  }

  return (
    <div className='container grid h-svh flex-col items-center justify-center bg-primary-foreground lg:max-w-none lg:px-0'>
      <div className='mx-auto flex w-full flex-col justify-center space-y-2 sm:w-[480px] lg:p-8'>
        <div className='mb-4 flex items-center justify-center'>
          <IconMessages size={20} className='mr-2 h-6 w-6' />
          <h1 className='text-xl font-medium'>Karlota Messenger</h1>
        </div>
        {children}
      </div>
    </div>
  )
}