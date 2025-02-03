import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { IconEdit, IconSearch } from "@tabler/icons-react";
import { Users } from "lucide-react";
import { useContact } from "@/features/chats/hooks/use-contact";
import { Suspense, useState, useEffect } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { NewGroupContactSelection } from "./new-group-contact-selection";
import { NewGroupCreation } from "./new-group-creation";
import { Profile } from "@/types/user";
import { NewChatCreation } from "@/features/chats/components/new-chat-creation";
import { useChatStore } from "@/states/chat-store";
import { useWebSocket } from "@/features/chats/context/websocket";
import {toast} from "@/hooks/use-toast";
import {ChatListState} from "@/features/chats/components/chat-list";
import {useGlobalActionStore} from "@/states/global-action-store";

enum ChatStage {
  Default,
  NewChat,
  ContactSelection,
  GroupCreation,
}

export function NewChatDropdown() {
  const { data: contacts } = useContact();
  const { chats, setSelectedChat } = useChatStore();
  const { sendMessage } = useWebSocket();
  const [search, setSearch] = useState("");
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [chatStage, setChatStage] = useState<ChatStage>(ChatStage.Default);
  const [selectedContact, setSelectedContact] = useState<Profile | null>(null);
  const [selectedContacts, setSelectedContacts] = useState<Profile[]>([]);
  const [newChatId, setNewChatId] = useState<number | null>(null);
  const { setState } = useGlobalActionStore();

  const filteredContactList = (contacts?.data ?? []).filter(({ display_name }) =>
    display_name?.toLowerCase().includes(search.trim().toLowerCase()));

  useEffect(() => {
    if (!newChatId) return;
    const chat = chats.find(
      (chat) => chat.id === newChatId);
    if (chat) {
      cleanupState();
      sendMessage({ action: "messages", conversation_id: chat.id });
      setSelectedChat(chat);
    }
  }, [chats, newChatId]);

  const onContactSelected = (user: Profile) => {
    const chat = chats.find((chat) =>
        chat.type === "private" && chat.users.some((u) => u.id === user.id));
    if (chat) {
      sendMessage({ action: "messages", conversation_id: chat.id });
      cleanupState();
      return setSelectedChat(chat);
    }
    setSelectedContact(user);
    setChatStage(ChatStage.NewChat);
  };

  const onNewPrivateChatCreated = (chatId: number) => {
    if (chatId === 0) {
      setChatStage(ChatStage.Default);
      return;
    }
    setNewChatId(chatId);
    setState(ChatListState, true)
    sendMessage({ action: "chats" });
  };

  const onNextGroupCreation = (selectedContacts: Profile[]) => {
    if (selectedContacts.length < 2) {
      toast({
        variant: 'destructive',
        title: 'Select at least 2 contacts!',
      });
      return;
    }
    setChatStage(ChatStage.GroupCreation);
    setSelectedContacts(selectedContacts);
  }

  const onNewGroupChatCreated = (chatId: number) => {
    setNewChatId(chatId);
    setState(ChatListState, true)
    sendMessage({ action: "chats" });
  }

  const cleanupState = () => {
    setSearch("");
    setChatStage(ChatStage.Default);
    setSelectedContact(null);
    setSelectedContacts([]);
    setNewChatId(null);
    setIsDropdownOpen(false);
  }

  return (
    <DropdownMenu open={isDropdownOpen} onOpenChange={setIsDropdownOpen}>
      <DropdownMenuTrigger asChild>
        <Button size="icon" variant="ghost" className="rounded-lg">
          <IconEdit size={24} className="stroke-muted-foreground" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="p-4 w-96 bg-gray-100">
        {chatStage === ChatStage.Default && (
          <>
            <DropdownMenuLabel className="text-center">New Chat</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <div className="space-y-4 mt-4">
              <label className="bg-white flex h-10 w-full items-center rounded-md border pl-2">
                <IconSearch size={15} className="mr-2 stroke-slate-500" />
                <input
                  type="text"
                  className="w-full flex-1 bg-inherit text-sm focus-visible:outline-none"
                  placeholder="Search"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </label>
              <Button
                variant="secondary"
                className="w-full h-10 bg-white"
                onClick={() => {
                  setSearch("")
                  setChatStage(ChatStage.ContactSelection)
                }}
              >
                <Users className="mr-2" /> New Group
              </Button>
            </div>
            <DropdownMenuLabel className="mt-4">Contacts on Karlota</DropdownMenuLabel>
            <Suspense fallback={<div className="text-center">Loading contacts...</div>}>
              <ScrollArea className="h-[400px] w-full p-2">
                <div className="space-y-4">
                  {filteredContactList.map((user) => (
                    <div
                      key={user.id}
                      className="flex items-center gap-4 rounded-lg border p-4 bg-white cursor-pointer hover:bg-gray-50"
                      onClick={() => onContactSelected(user)}
                    >
                      <Avatar>
                        <AvatarFallback>
                          {user.display_name
                            .split(" ")
                            .map((n) => n[0])
                            .join("")}
                        </AvatarFallback>
                      </Avatar>
                      <div className="text-left">
                        <p className="text-sm font-medium">{user.display_name}</p>
                        <p className="text-sm text-muted-foreground">{user.email}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </Suspense>
          </>
        )}
        {chatStage === ChatStage.NewChat && (
          <NewChatCreation
            onBack={onNewPrivateChatCreated}
            selectedContact={selectedContact}
          />
        )}
        {chatStage === ChatStage.ContactSelection && (
          <NewGroupContactSelection
            contacts={filteredContactList}
            onBack={() => {
              setChatStage(ChatStage.Default);
              setSelectedContacts([]);
            }}
            onNext={onNextGroupCreation}
            selectedContact={selectedContacts}
          />
        )}
        {chatStage === ChatStage.GroupCreation && (
          <NewGroupCreation
            onBack={() => setChatStage(ChatStage.ContactSelection)}
            totalContacts={contacts?.data?.length as number}
            selectedContacts={selectedContacts}
            setSelectedContacts={setSelectedContacts}
            onCreate={onNewGroupChatCreated}
          />
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}