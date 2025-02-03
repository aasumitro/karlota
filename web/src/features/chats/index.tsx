import {useState} from 'react'
import {
  IconMessages,
  IconSearch,
} from '@tabler/icons-react'
import { cn } from '@/lib/utils'
import {UserActionDropdown} from "@/components/user-action-dropdown";
import {useNavigate} from "@tanstack/react-router";
import {useAuthStore} from "@/states/auth-store";
import {toast} from "@/hooks/use-toast";
import {isJWTExpired} from "@/lib/jwt";
import {LogoutModal} from "@/components/logout-modal";
import {ProfileInfoModal} from "@/components/profile-info-modal";
import {UpdateDisplayNameModal} from "@/components/update-displayname-modal";
import {UpdatePasswordModal} from "@/components/update-password-modal";
import {NewChatDropdown} from "@/features/chats/components/new-chat-dropdown";
import {ChatNone} from "@/features/chats/components/chat-none";
import {ChatList} from "@/features/chats/components/chat-list";
import {ChatDetail} from "@/features/chats/components/chat-detail";
import {useChatStore} from "@/states/chat-store";
import {WebSocketProvider} from "@/features/chats/context/websocket";
import {LeaveGroupModal} from "@/features/chats/components/leave-group-modal";
import {DeleteGroupModal} from "@/features/chats/components/delete-group-modal";

export default function Chats() {
  const navigate = useNavigate()
  if (!useAuthStore.getState().auth.accessToken &&
    isJWTExpired(useAuthStore.getState().auth.accessToken)) {
    navigate({to: "/sign-in", replace: true}).then(() =>
      toast({ variant: 'destructive', title: "token expired" }))
  }

  const [search, setSearch] = useState('')
  const {chats, selectedChat, mobileSelectedChat} = useChatStore();

  // Filtered data based on the search query
  const filteredChatList = chats.filter(({ type, name, users}) => {
    const chatName =
      type === "private"
        ? users.map((user) => user.display_name).join(", ")
        : name;
    return chatName.toLowerCase().includes(search.trim().toLowerCase());
  });

  return (
    <main className='peer-[.header-fixed]/header:mt-16 fixed-main flex flex-col flex-grow overflow-hidden'>
      <WebSocketProvider>
        <section className='flex h-full'>
          {/* Left Side */}
          <div className='flex w-full flex-col sm:w-56 lg:w-72 2xl:w-80'>
            <div className='sticky top-0 z-10 bg-background p-4 pb-3 shadow-md sm:static sm:z-auto sm:shadow-none'>
              <div className='flex items-center justify-between py-2 mb-4'>
                <div className='flex gap-2'>
                  <h1 className='text-2xl font-bold'>Karlota</h1>
                  <IconMessages size={20} />
                </div>
                <NewChatDropdown />
              </div>

              <label className='flex h-12 w-full items-center space-x-0 rounded-md border border-input pl-2 focus-within:outline-none focus-within:ring-1 focus-within:ring-ring'>
                <IconSearch size={15} className='mr-2 stroke-slate-500' />
                <span className='sr-only'>Search</span>
                <input
                  type='text'
                  className='w-full flex-1 bg-inherit text-sm focus-visible:outline-none'
                  placeholder='Search chat...'
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </label>
            </div>

            <ChatList chats={filteredChatList}/>

            <div className='sticky bottom-0 z-10 bg-background shadow-md sm:static sm:z-auto sm:shadow-none'>
              <UserActionDropdown />
            </div>
          </div>

          {/* Right Side */}
          <div className={cn(
            'absolute inset-0 hidden left-full z-50 w-full flex-1 flex-col rounded-md ',
            'border bg-primary-foreground shadow-sm transition-all duration-200 sm:static sm:z-auto sm:flex ',
            mobileSelectedChat && 'left-0 flex'
          )}>
            {selectedChat ? <ChatDetail/> : <ChatNone />}
          </div>
        </section>
        <DeleteGroupModal />
        <LeaveGroupModal/>
      </WebSocketProvider>

      <ProfileInfoModal />
      <UpdateDisplayNameModal />
      <UpdatePasswordModal />
      <LogoutModal/>
    </main>
  )
}