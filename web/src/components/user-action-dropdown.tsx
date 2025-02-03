import {Avatar, AvatarFallback} from "@/components/ui/avatar";
import {useState} from "react";
import {ChevronsDownUp, ChevronsUpDown, LogOut, RectangleEllipsis, Settings2, User2, UserRoundPen} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuPortal, DropdownMenuItem,
} from "@/components/ui/dropdown-menu"
import {useAuthStore} from "@/states/auth-store";
import {useProfile} from "@/hooks/use-account";
import {isJWTExpired} from "@/lib/jwt";
import {ProfileInfoModalState} from "@/components/profile-info-modal";
import {UpdateDisplayNameModalState} from "@/components/update-displayname-modal";
import {UpdatePasswordModalState} from "@/components/update-password-modal";
import {useGlobalActionStore} from "@/states/global-action-store";
import {LogoutModalState} from "@/components/logout-modal";
import {useSelectActionStore} from "@/states/select-action-store";

export const UserActionDropdown = () => {
  const [openMenu, setOpenMenu] = useState(false)
  const { auth } = useAuthStore();
  const { setState } = useGlobalActionStore();
  const { setUser } = useSelectActionStore();

  if (!auth.user && !isJWTExpired(auth.accessToken)) {
    const { data: profile, isPending } = useProfile();
    auth.setUser(profile.data);
    if (!auth.user && isPending) {
      return <>Loading . . .</>;
    }
  }

  const showUserInfo = () => {
    setState(ProfileInfoModalState, true);
    setUser(auth.user);
  }

  return (
    <DropdownMenu open={openMenu} onOpenChange={setOpenMenu}>
      <DropdownMenuTrigger asChild>
        <button className="py-5 px-4 flex justify-between  bg-gray-100 hover:bg-gray-200 dark:hover:bg-gray-800 items-center w-full">
          <div className={`flex items-center gap-2.5`}>
            <Avatar className="w-10 h-10 rounded-full flex items-center justify-center bg-purple-100">
              <AvatarFallback className="text-2xl">
                {auth.user?.display_name[0]?.toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className={`text-left`}>
              <p className="font-semibold text-xs mb-0.5">
                {auth.user?.display_name}
              </p>
              <p className="font-medium text-xs text-gray-500">
                {auth.user?.email}
              </p>
            </div>
          </div>
          <div
            className={`transition-transform duration-300 ease-in-out ${
              openMenu ? "rotate-180" : "rotate-0"
            }`}
          >
            {openMenu ? <ChevronsDownUp size={18} /> : <ChevronsUpDown size={18} />}
          </div>
        </button>
      </DropdownMenuTrigger>
      {/* make validation if web right if mobile top*/}
      <DropdownMenuContent className="w-56" side= "top">
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onClick={showUserInfo}
          className="flex w-full"
        >
          <User2 className="mr-2 h-4 w-4" />
          <span>Profile</span>
        </DropdownMenuItem>
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>
            <Settings2 className="mr-2 h-4 w-4" />
            <span>Update</span>
          </DropdownMenuSubTrigger>
          <DropdownMenuPortal>
            <DropdownMenuSubContent>
              <DropdownMenuItem
                onClick={() => setState(UpdateDisplayNameModalState, true)}
                className="flex w-full"
              >
                <UserRoundPen className="mr-2 h-4 w-4"  />
                <span>Display Name</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => setState(UpdatePasswordModalState, true)}
                className="flex w-full"
              >
                <RectangleEllipsis className="mr-2 h-4 w-4"  />
                <span>Password</span>
              </DropdownMenuItem>
            </DropdownMenuSubContent>
          </DropdownMenuPortal>
        </DropdownMenuSub>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onClick={() => setState(LogoutModalState, true)}
          className="flex w-full"
        >
          <LogOut className="mr-2 h-4 w-4" />
          <span>Log out</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}