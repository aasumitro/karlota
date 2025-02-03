import { Button } from "@/components/ui/button"
import { ArrowLeft,Check } from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import {Profile} from "@/types/user";
import {ScrollArea} from "@/components/ui/scroll-area";
import {useState} from "react";
import {IconSearch} from "@tabler/icons-react";

interface ContactSelectionProps {
  contacts?: Profile[]
  selectedContact: Profile[]
  onBack: () => void
  onNext?: (selectedContacts: Profile[]) => void
}

export function NewGroupContactSelection({ contacts, selectedContact, onBack, onNext }: ContactSelectionProps) {
  const [selectedContacts, setSelectedContacts] = useState<Profile[]>(selectedContact)
  const [search, setSearch] = useState("")

  const filteredContactList = contacts?.filter(({ display_name }) =>
    display_name.toLowerCase().includes(search.trim().toLowerCase()))

  const toggleContact = (contact: Profile) => {
    setSelectedContacts((prev) => {
      const isSelected = prev.some((c) => c.id === contact.id)
      if (isSelected) {
        return prev.filter((c) => c.id !== contact.id)
      } else {
        return [...prev, contact]
      }
    })
  }

  const handleNext = () => {
    if (onNext && selectedContacts.length > 0) {
      onNext(selectedContacts)
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-4">
        <Button variant="ghost" size="sm" onClick={onBack}>
          <ArrowLeft className="h-4 w-4 mr-2" />
        </Button>
        <div className="text-sm font-semibold flex flex-col items-center">
          Add Members
          <span className="text-xs font-light">
            {selectedContacts.length}/{contacts?.length}
          </span>
        </div>
        <Button
          variant="ghost" size="sm"
          disabled={selectedContacts.length === 0}
          onClick={handleNext}
        >
          Next
        </Button>
      </div>

      <div className="space-y-4 mt-4">
        <label className="bg-white flex h-10 w-full items-center space-x-0 rounded-md border border-input pl-2 focus-within:outline-none focus-within:ring-1 focus-within:ring-ring">
          <IconSearch size={15} className="mr-2 stroke-slate-500" />
          <span className="sr-only">Search</span>
          <input
            type="text"
            className="w-full flex-1 bg-inherit text-sm focus-visible:outline-none"
            placeholder="Search"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </label>
      </div>

      <ScrollArea className="h-[400px] w-full">
        <div className="space-y-4">
          {filteredContactList?.map((user) => {
            const isSelected = selectedContacts.some((c) => c.id === user.id)
            return (
              <div
                key={user.id}
                className={`flex items-center justify-between space-x-4 rounded-lg border p-4 transition-colors bg-white cursor-pointer hover:bg-gray-50 ${
                  isSelected ? "border-primary" : ""
                }`}
                onClick={() => toggleContact(user)}
              >
                <div className="flex items-center space-x-4">
                  <Avatar>
                    <AvatarFallback>
                      {user.display_name
                        .split(" ")
                        .map((n) => n[0])
                        .join("")}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="text-sm font-medium leading-none">{user.display_name}</p>
                    <p className="text-sm text-muted-foreground">{user.email}</p>
                  </div>
                </div>
                {isSelected && (
                  <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary">
                    <Check className="w-4 h-4 text-primary-foreground" />
                  </div>
                )}
              </div>
            )
          })}
        </div>
      </ScrollArea>
    </div>
  )
}

