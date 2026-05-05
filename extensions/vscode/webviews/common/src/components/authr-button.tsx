import { useNavigate, Link } from "@tanstack/react-router";
import { useQueryClient, useMutation } from '@tanstack/react-query'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  // DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from "./ui/dropdown-menu"
import { useAuthr } from "../provider/authr";

export const AuthrButton = ({ children }: { children?: React.ReactNode }) => {
  
  const navigate = useNavigate();
  const authr = useAuthr();
  const queryClient = useQueryClient()

  const switchAccount = useMutation({
    mutationFn: async (did: string) => {
      authr.switchAccount(did)
      return null
    },
    onSuccess: (_, did: string) => {
      console.log("invalidating acct queries", did)
      queryClient.invalidateQueries({ queryKey: [did, 'acct'] })
    },
  })

  if (authr?.session?.did) {
    return (
      <DropdownMenu>
        <DropdownMenuTrigger
          className="bg-white text-black px-2 py-1 rounded hover:bg-blue-200"
        >@{authr.session.handle}</DropdownMenuTrigger>
        <DropdownMenuContent
          className="w-56 bg-white text-black"
        >
          {/* <DropdownMenuLabel>
            hello
          </DropdownMenuLabel>
          <DropdownMenuSeparator /> */}

          <DropdownMenuItem onSelect={() => navigate({ to: "/account" })}>Account</DropdownMenuItem>

          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Switch account...</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent
                className="bg-white text-black"
              >
                { authr?.sessions?.accounts?.map((s: any) => {
                  return (
                    <DropdownMenuItem key={s.did} onSelect={() => switchAccount.mutate(s.did) }>
                      {s.handle}
                    </DropdownMenuItem>
                  )
                })}
                <DropdownMenuSeparator />
                <DropdownMenuItem onSelect={() => navigate({ to: "/sign-in" })}>Add account</DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>

          { children ? <DropdownMenuSeparator /> : null }

          { children }

        </DropdownMenuContent>
      </DropdownMenu>
    );
  }

  return (
    <Link to="/sign-in" className="bg-white text-black px-2 py-1 rounded hover:bg-blue-200">
      Sign In
    </Link>
  );
}