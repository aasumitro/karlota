import {Fragment} from "react/jsx-runtime";
import {Skeleton} from "@/components/ui/skeleton";
import {Separator} from "@/components/ui/separator";

interface props {
  len: number;
}

export function ChatListSkeleton({ len }: props) {
  return (
    Array.from({ length: len }).map((_, index) => (
      <Fragment key={index}>
        <button type="button" className="flex w-full rounded-md px-2 py-2 text-left text-sm hover:bg-secondary/75">
          <div className="flex gap-2 items-center w-full">
            {/* Avatar Wrapper to Prevent Stretching */}
            <div className="h-10 w-10 shrink-0">
              <Skeleton className="h-10 w-10 rounded-full aspect-square" />
            </div>
            {/* Name & Message Skeletons */}
            <div className="flex flex-col flex-1 space-y-1">
              <Skeleton className="h-4 w-32" />
              <Skeleton className="h-4 w-40" />
            </div>
          </div>
        </button>
        <Separator className="my-1" />
      </Fragment>
    ))
  )
}