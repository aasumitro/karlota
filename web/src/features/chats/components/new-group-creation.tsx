import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ArrowLeft, Camera } from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { useState } from "react"
import {Profile} from "@/types/user";
import {useConversation} from "@/features/chats/hooks/use-conversation";
import {toast} from "@/hooks/use-toast";

interface GroupCreationProps {
  onBack: () => void
  totalContacts: number
  selectedContacts: Profile[]
  setSelectedContacts: (contacts: Profile[]) => void
  onCreate: (chatId: number) => void
}

export function NewGroupCreation({ onBack, totalContacts, selectedContacts, setSelectedContacts, onCreate }: GroupCreationProps) {
  const [groupName, setGroupName] = useState("")
  const { mutate: newGroup, isPending} = useConversation();

  const handleCreate = () => {
    const members = selectedContacts.map(contact => ({
      id: contact.id,
      display_name: contact.display_name,
    }));
    newGroup({
      name: groupName,
      members: members,
    }, {
      onSuccess: async (response) => {
        const resp = response?.data
        if (!resp) {
          return;
        }
        if (onCreate) {
          onCreate(Number(resp))
        }
      }
    })
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-4">
        <Button variant="ghost" size="sm" onClick={onBack}>
          <ArrowLeft className="h-4 w-4 mr-2" />
        </Button>
        <div className="text-sm font-semibold">New group</div>
        <Button
          variant="ghost" size="sm"
          disabled={selectedContacts.length === 0 || isPending}
          onClick={handleCreate}
        >
          {isPending ? "Creating..." : "Create"}
        </Button>
      </div>

      <div className="flex items-center space-x-4 p-4 bg-white rounded-lg">
        <div className="relative">
          <Avatar className="h-16 w-16">
            <AvatarFallback>NG</AvatarFallback>
          </Avatar>
          <Button
            size="icon"
            variant="secondary"
            className="absolute bottom-0 right-0 h-6 w-6 cursor-pointer"
            onClick={() => {
              toast({
                variant: 'destructive',
                title: 'Currently not supported!',
              })
            }}
          >
            <Camera className="h-4 w-4" />
          </Button>
        </div>
        <Input
          placeholder="Group name (optional)"
          value={groupName}
          onChange={(e) => setGroupName(e.target.value)}
          className="flex-1"
        />
      </div>

      <div className="pt-4">
        <div className="text-sm text-muted-foreground mb-2">
          MEMBERS: {selectedContacts.length} OF {totalContacts}
        </div>
        <ScrollArea className="h-[200px]">
          <div className="space-y-2">
            {selectedContacts.map((contact) => (
              <div key={contact.id} className="flex items-center justify-between p-2 bg-white rounded-lg">
                <div className="flex items-center space-x-2">
                  <Avatar>
                    <AvatarFallback>
                      {contact.display_name
                        .split(" ")
                        .map((n: string) => n[0])
                        .join("")}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="text-sm font-medium leading-none">{contact.display_name}</p>
                    <p className="text-sm text-muted-foreground">{contact.email}</p>
                  </div>
                </div>
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => setSelectedContacts(selectedContacts.filter((c) => c.id !== contact.id))}
                >
                  Remove
                </Button>
              </div>
            ))}
          </div>
        </ScrollArea>
      </div>
    </div>
  )
}

