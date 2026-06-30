"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { Logo } from "@/components/icons";

export default function RegisterPage() {
  const router = useRouter();

  // Mock-only handler: nothing is sent anywhere. We simply prevent the
  // default form navigation and route the visitor to the demo account page.
  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    router.push("/account");
  }

  return (
    <section className="bg-white">
      <div className="flex min-h-[60vh] items-center justify-center px-4 py-12">
        <div className="w-full max-w-md rounded-lg border border-botb-card-border bg-white p-8 shadow-sm">
          <div className="flex flex-col items-center text-center">
            <Logo className="h-10 w-auto" />
            <h1 className="mt-6 font-jost text-[28px] font-bold uppercase text-botb-text">
              Sign up
            </h1>
            <p className="mt-2 text-[13px] text-botb-muted">
              Demo only — this is a portfolio clone; no account is created and no
              data is sent.
            </p>
          </div>

          <form onSubmit={handleSubmit} className="mt-8 flex flex-col gap-4">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="flex flex-col gap-1.5">
                <label
                  htmlFor="firstName"
                  className="text-[14px] font-medium text-botb-text"
                >
                  First name
                </label>
                <input
                  id="firstName"
                  name="firstName"
                  type="text"
                  required
                  autoComplete="given-name"
                  className="w-full rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
                />
              </div>

              <div className="flex flex-col gap-1.5">
                <label
                  htmlFor="lastName"
                  className="text-[14px] font-medium text-botb-text"
                >
                  Last name
                </label>
                <input
                  id="lastName"
                  name="lastName"
                  type="text"
                  required
                  autoComplete="family-name"
                  className="w-full rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
                />
              </div>
            </div>

            <div className="flex flex-col gap-1.5">
              <label
                htmlFor="email"
                className="text-[14px] font-medium text-botb-text"
              >
                Email
              </label>
              <input
                id="email"
                name="email"
                type="email"
                required
                autoComplete="email"
                className="w-full rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label
                htmlFor="password"
                className="text-[14px] font-medium text-botb-text"
              >
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                required
                autoComplete="new-password"
                className="w-full rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
              />
            </div>

            <div className="flex flex-col gap-2 text-[14px] text-botb-muted">
              <label className="flex items-start gap-2">
                <input
                  type="checkbox"
                  name="terms"
                  required
                  className="mt-0.5 h-4 w-4 accent-botb-orange"
                />
                <span>
                  I agree to the{" "}
                  <Link
                    href="/terms"
                    className="text-botb-orange hover:text-botb-orange-hover"
                  >
                    Terms of Play
                  </Link>
                </span>
              </label>
              <label className="flex items-start gap-2">
                <input
                  type="checkbox"
                  name="ageConfirm"
                  required
                  className="mt-0.5 h-4 w-4 accent-botb-orange"
                />
                <span>I am 18 or over</span>
              </label>
            </div>

            <button
              type="submit"
              className="w-full rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
            >
              Sign up
            </button>
          </form>

          <p className="mt-6 text-center text-[14px] text-botb-muted">
            Already have an account?{" "}
            <Link
              href="/login"
              className="font-medium text-botb-orange hover:text-botb-orange-hover"
            >
              Log in
            </Link>
          </p>
        </div>
      </div>
    </section>
  );
}
