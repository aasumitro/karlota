import {Profile} from "@/types/user";
import {Button} from "@/components/ui/button";
import {ArrowLeft} from "lucide-react";
import {useState} from "react";
import {useConversation} from "@/features/chats/hooks/use-conversation";
import {toast} from "@/hooks/use-toast";
import {Textarea} from "@/components/ui/textarea";

interface ChatCreationProps {
  onBack: (chatId: number) => void
  selectedContact: Profile | null
}

export function NewChatCreation({ onBack, selectedContact }: ChatCreationProps) {
  const [message, setMessage] = useState("")
  const { mutate: newChat, isPending} = useConversation();

  const handleCreate = () => {
    if (!message) {
      toast({
        variant: 'destructive',
        title: 'Message is required!',
      })
      return;
    }
    newChat({
      message: message,
      members: [{
        id: selectedContact?.id,
        display_name: selectedContact?.display_name
      }],
    }, {
      onSuccess: async (response) => {
        const resp = response?.data
        if (!resp) {
          return;
        }
        if (onBack) {
          onBack(Number(resp));
        }
      }
    })
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-4">
        <Button variant="ghost" size="sm" onClick={() => {onBack(0)}}>
          <ArrowLeft className="h-4 w-4 mr-2"/>
        </Button>
        <div className="text-sm font-semibold">
          Say hello to {selectedContact?.display_name}!
        </div>
        <Button
          variant="ghost" size="sm"
          disabled={message.length === 0 || isPending}
          onClick={handleCreate}
        >
          {isPending ? "Sending..." : "Sent"}
        </Button>
      </div>

      <div className="flex items-center space-x-4 p-2 bg-white rounded-lg">
        <Textarea
          placeholder={`Hi, how are you ${selectedContact?.display_name}?`}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          className="flex-1"
        />
      </div>
    </div>
  )
}
