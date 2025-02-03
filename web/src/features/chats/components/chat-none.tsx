import {IconMessages} from "@tabler/icons-react";

export function ChatNone() {
  return (
    <div className='m-auto flex h-full w-full flex-col items-center justify-center'>
      <IconMessages size={100} />
      <h1 className='text-[1rem] font-semibold tracking-wide mt-8'>
        Karlota Messenger for Web
      </h1>
    </div>
  )
}