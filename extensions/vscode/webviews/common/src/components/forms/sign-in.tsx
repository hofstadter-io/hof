"use client"

import { useSearch } from "@tanstack/react-router"
import { useForm } from "@tanstack/react-form"
import { z } from "zod"

import { cn } from "../../lib/utils"
import { useAuthr } from "../../provider/authr"

const formSchema = z.object({
  handle: z.string().min(2, {
    message: "Handle must be at least 2 characters.",
  }),
  redirect: z.string().optional(),
})

export function AuthrSignInForm(props: any) {
  const searchParams: any = useSearch({
    strict: false
  });
  const { login } = useAuthr()

  const form = useForm({
    defaultValues: {
      handle: searchParams.handle || "",
    },
    validators: {
      onChange: formSchema,
    },

    onSubmit: async ({ value }: { value: any }) => {
      const baseUrl = window.location.origin
      if (value.redirect) {
        if (value.redirect.startsWith("/")) {
          // this is a relative path, so we need to make it absolute
          value.redirect = baseUrl + value.redirect
        }
      } else {
        value.redirect = baseUrl + (props.redirect || "/")
      }

      console.log(value)

      login(value)
    }
  })

  return (
    <div>

      <div className="flex flex-col items-center px-2">
        <span className="flex items-center text-xl font-light">
          Login to Veggie
        </span>
        <span className="text-center">a healthy agent framework</span>
      </div>

      <form
        className="space-y-8 flex flex-col gap-1"
        onSubmit={(e) => {
          // console.log("form.onSubmit", e)
          e.preventDefault()
          e.stopPropagation()
        }}
      >

        <form.Field
          // control={form.control}
          name="handle"
          children={(field) => (
            <div className="flex flex-col">
              <label
                className="m-1 text-md font-light"
                htmlFor={field.name}
              >{field.name}</label>
              <input
                name={field.name}
                className={cn(
                  "border rounded p-2 m-0",
                  !(field.state.meta.isPristine || field.state.meta.isValid) ? "border-red-500" : "border-gray-300",
                )}
                value={field.state.value || ''}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
                placeholder="username.bsky.social"
              />
              {
                field.state.meta.isBlurred && !field.state.meta.isValid &&
                <em
                  className="m-1 text-red-500"
                  role="alert"
                >{field.state.meta.errors.map((e) => e?.message).join(', ')}</em>
              }
              {/* <pre>{JSON.stringify(field.state.meta, null, 2)}</pre> */}
            </div>
          )}
        />


        <div className="flex flex-row gap-4">
          <button
            type="submit"
            className="mx-2 px-3 py-1 bg-blue-500 rounded text-white"
            onClick={() => form.handleSubmit({ submitAction: 'post' })}
          >Sign in</button>
        </div>
      </form>
    </div>
  )
}

