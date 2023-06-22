import { authConfig } from '@/src/lib/auth/config'
import NextAuth from "next-auth"

const handler = NextAuth(authConfig)

export { handler as GET, handler as POST }
