import {Button} from "@/components/ui/button";
import {IconArrowLeft, IconDotsVertical, IconPhone, IconVideo} from "@tabler/icons-react";
import {Avatar, AvatarFallback} from "@/components/ui/avatar";
import {useChatStore} from "@/states/chat-store";
import {useEffect, useState} from "react";
import {useAuthStore} from "@/states/auth-store";
import {toast} from "@/hooks/use-toast";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {useGlobalActionStore} from "@/states/global-action-store";
import {LeaveGroupModalState} from "@/features/chats/components/leave-group-modal";
import {cn} from "@/lib/utils";
import {ProfileInfoModalState} from "@/components/profile-info-modal";
import {useSelectActionStore} from "@/states/select-action-store";
import {Profile} from "@/types/user";
import {DeleteGroupModalState} from "@/features/chats/components/delete-group-modal";
import {NewCallModalState, NewCallState} from "@/features/chats/components/new-call";
import {CallStage} from "@/types/chat";

export function ChatToolbar() {
  const sessionID = useAuthStore.getState().auth.user?.id;
  const { selectedChat, setSelectedChat, setMessages } = useChatStore();
  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [currentUser, setCurrentUser] = useState<Profile | null>(null)
  const [otherUser, setOtherUser] = useState<Profile | null>(null)
  const {setState, setStatus} = useGlobalActionStore();
  const { setUser } = useSelectActionStore();

  const notSupportedToast = () => {
    toast({
      variant: 'destructive',
      title: 'Currently not supported!',
    })
  }

  useEffect(() => {
    setDescription("")
    if (selectedChat) {
      if (selectedChat.type === "private") {
        const otherUser = selectedChat.users.find((user) => user.id !== sessionID);
        setName(otherUser?.display_name || "Unknown");
        if (otherUser) {
          setOtherUser(otherUser as Profile)
          setDescription("Click here to see info")
        }
      } else if (selectedChat.type === "group") {
        setName(selectedChat?.name || "Group Chat");
        let desc = `${selectedChat.users.length} Members, `
        desc += `${selectedChat.users.filter((u) => u.is_online).length} Online`
        const currentUser = selectedChat.users.find((user) => user.id === sessionID);
        if (currentUser) {
          setCurrentUser(currentUser)
        }
        setDescription(desc)
      }
    }
  }, [selectedChat, sessionID]);

  const closeComponent = () => {
    setSelectedChat(null)
    setMessages(null)
    setName("")
    setDescription("")
    setOtherUser(null)
    setCurrentUser(null)
  }

  return (
    <div className='mb-1 flex flex-none justify-between rounded-t-md bg-secondary p-4 shadow-lg'>
      {/* Left */}
      <div className='flex gap-3'>
        <Button
          size='icon'
          variant='ghost'
          className='-ml-2 h-full sm:hidden'
          onClick={closeComponent}
        >
          <IconArrowLeft />
        </Button>
        <div className={cn(
          'flex items-center gap-2 lg:gap-4',
         (selectedChat?.type === "private") && "cursor-pointer"
        )} onClick={() => {
          if (selectedChat?.type === "private") {
            setState(ProfileInfoModalState, true);
            setUser(otherUser);
          }
        }}>
          <Avatar className='size-9 lg:size-11'>
            <AvatarFallback className="bg-white">{name[0]?.toUpperCase() ?? "U"}</AvatarFallback>
          </Avatar>
          <div>
            <span className='col-start-2 row-span-2 text-sm font-medium lg:text-base tracking-tighter hover:border-b-2 hover:border-black'>
              {name}
            </span>
            <span className='col-start-2 row-span-2 row-start-2 line-clamp-1 block max-w-32 text-ellipsis text-nowrap text-xs text-muted-foreground lg:max-w-none lg:text-sm'>
              {otherUser?.is_online ? "online": description} 
            </span>
          </div>
        </div>
      </div>

      {/* Right */}
      <div className='-mr-1 flex items-center gap-1 lg:gap-2'>
        {selectedChat?.type === "private" && <Button
          size='icon'
          variant='ghost'
          className='hidden size-8 rounded-full sm:inline-flex lg:size-10'
          onClick={() => {
            setState(NewCallModalState, true)
            setStatus(NewCallState, CallStage.Calling)
          }}
        >
          <IconVideo size={22} className='stroke-muted-foreground' />
        </Button>}
        {selectedChat?.type === "private" &&  <Button
          size='icon'
          variant='ghost'
          className='hidden size-8 rounded-full sm:inline-flex lg:size-10'
          disabled
          // onClick={() =>  {
          //   setState(NewCallModalState, true)
          //   setStatus(NewCallState, CallStage.Calling)
          // }}
        >
          <IconPhone size={22} className='stroke-muted-foreground' />
        </Button>}
        <DropdownMenu>
          <DropdownMenuTrigger>
            <IconDotsVertical className='stroke-muted-foreground sm:size-5' />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {selectedChat?.type === "group" && <DropdownMenuItem onClick={notSupportedToast}>
              Detail
            </DropdownMenuItem>}
            {selectedChat?.type === "private" && <>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={notSupportedToast}
                className="text-red-500"
              >Delete chat</DropdownMenuItem>
            </>}
            {selectedChat?.type === "group" && currentUser?.role === "admin" && <>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-red-500"
                onClick={() => setState(DeleteGroupModalState, true)}
              >Delete group</DropdownMenuItem>
            </>}
            {selectedChat?.type === "group" && currentUser?.role === "member"  && <>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-red-500"
                onClick={() => setState(LeaveGroupModalState, true)}
              >Exit group</DropdownMenuItem>
            </>}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}