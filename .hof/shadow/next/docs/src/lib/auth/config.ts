import type { Adapter } from "next-auth/adapters";
import { NextAuthOptions, User, getServerSession } from "next-auth";
import { useSession } from "next-auth/react";
import { redirect } from "next/navigation";

import CredentialsProvider from "next-auth/providers/credentials";
import GoogleProvider from "next-auth/providers/google";
// import GithubProvider from "next-auth/providers/github";

import { PrismaAdapter } from "@auth/prisma-adapter";
import { db } from "../db";

export const authConfig: NextAuthOptions = {
  adapter: PrismaAdapter(db) as Adapter,
  session: {
    maxAge: 30 * 24 * 60 * 60,
  },
  providers: [
    CredentialsProvider({
      name: "Email",
      credentials: {
        email: {
          label: "Email",
          type: "email",
          placeholder: "user@domain.com",
        },
        password: {
          label: "Password",
          type: "password",
        },
      },
      async authorize(credentials) {
        if (!credentials || !credentials.email || !credentials.password) {
          return null;
        }

        const dbUser = await db.user.findFirst({
          where: { email: credentials.email },
        });
        console.log(credentials.email, dbUser);

        //// bcrypt credentials.password and compare, there should be a function for this
        //if (dbUser && dbUser.password === credentials.password) {
        //  // interesting pattern...
        //  const { password, ...cleanedUser } = dbUser;
        //  return cleanedUser as User;
        //}

        return null;
      },
    }),
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID! as string,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET! as string,
    }),
    //GithubProvider({
    //  clientId: process.env.GITHUB_CLIENT_ID as string,
    //  clientSecret: process.env.GITHUB_CLIENT_SECRET as string,
    //}),
  ],
  callbacks: {
    async session({ token, session }) {
      if (token) {
        session.user.id = token.id;
        session.user.name = token.name;
        session.user.email = token.email;
        session.user.image = token.image;
        session.user.username = token.username;
      }
      return session;
    },
    async jwt({ token, user }) {
      const dbUser = await db.user.findFirst({
        where: { email: token.email },
      });

      // if no existing user...
      if (!dbUser) {
        token.id = user!.id;
        return token;
      }

      // if user, but no username...
      if (!dbUser.username) {
        await db.user.update({
          where: {
            id: dbUser.id,
          },
          data: {
            username: dbUser.email,
          },
        });
      }

      return {
        id: dbUser.id,
        name: dbUser.name,
        email: dbUser.email,
        image: dbUser.image,
        username: dbUser.username,
      };
    },
    redirect() {
      return "/";
    },
  },
};
