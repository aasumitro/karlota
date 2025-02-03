import {IconMessages} from "@tabler/icons-react";
import {useEffect} from "react";
import {useAuthStore} from "@/states/auth-store";
import {useNavigate} from "@tanstack/react-router";

export default function Launcher() {
  const navigate = useNavigate();

  useEffect(() => {
    const timeoutId = setTimeout(async () => {
      if (!useAuthStore.getState().auth.accessToken) {
        await navigate({ to: "/sign-in", replace: true });
        return;
      }
      await navigate({ to: "/chats", replace: true });
    }, 1000);

    return () => clearTimeout(timeoutId);
  }, [])

  return (
    <div className='h-svh'>
      <div className='m-auto flex h-full w-full flex-col items-center justify-center gap-2'>
        <IconMessages size={150} />
        <h1 className='text-[3rem] font-semibold tracking-wide'>
          Karlota Messenger
        </h1>
      </div>
    </div>
  )
}