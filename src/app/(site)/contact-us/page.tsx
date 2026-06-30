"use client";

import { useState } from "react";
import Link from "next/link";

export default function ContactUsPage() {
  const [submitted, setSubmitted] = useState(false);

  // Local-only handler: this template has no backend, so we never submit
  // anywhere. We just prevent the default navigation and show a confirmation.
  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitted(true);
  }

  return (
    <section className="bg-white">
      <div className="mx-auto max-w-[1360px] px-4 py-10">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Contact Us
        </h1>
        <p className="mt-2 max-w-2xl text-botb-muted">
          Got a question about a competition, your account or a recent win? Our
          team is here to help. Drop us a message below or reach us directly
          using the details on the right.
        </p>

        <div className="mt-8 grid grid-cols-1 gap-10 lg:grid-cols-2">
          {/* Left: contact form */}
          <div>
            {submitted ? (
              <div className="rounded-md border border-botb-card-border bg-botb-gray p-6">
                <h2 className="font-jost text-[20px] font-semibold text-botb-text">
                  Thanks — we&apos;ll be in touch!
                </h2>
                <p className="mt-2 text-botb-muted">
                  Your message has been received. A member of the BOTB team will
                  get back to you as soon as possible.
                </p>
                <button
                  type="button"
                  onClick={() => setSubmitted(false)}
                  className="mt-4 rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
                >
                  Send another message
                </button>
              </div>
            ) : (
              <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                <div className="flex flex-col gap-1.5">
                  <label
                    htmlFor="name"
                    className="text-[14px] font-medium text-botb-text"
                  >
                    Name
                  </label>
                  <input
                    id="name"
                    name="name"
                    type="text"
                    required
                    autoComplete="name"
                    className="rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
                  />
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
                    className="rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <label
                    htmlFor="subject"
                    className="text-[14px] font-medium text-botb-text"
                  >
                    Subject
                  </label>
                  <input
                    id="subject"
                    name="subject"
                    type="text"
                    required
                    className="rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <label
                    htmlFor="message"
                    className="text-[14px] font-medium text-botb-text"
                  >
                    Message
                  </label>
                  <textarea
                    id="message"
                    name="message"
                    required
                    rows={6}
                    className="rounded-md border border-botb-card-border px-3 py-2 outline-none focus:border-botb-orange"
                  />
                </div>

                <div>
                  <button
                    type="submit"
                    className="rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white hover:bg-botb-orange-hover"
                  >
                    Send
                  </button>
                </div>
              </form>
            )}
          </div>

          {/* Right: contact details + FAQ */}
          <div className="flex flex-col gap-8">
            <div className="rounded-md border border-botb-card-border p-6">
              <h2 className="font-jost text-[20px] font-semibold text-botb-text">
                Get in touch
              </h2>
              <dl className="mt-4 flex flex-col gap-4 text-[15px]">
                <div>
                  <dt className="font-medium text-botb-text">Email</dt>
                  <dd className="text-botb-muted">
                    <a
                      href="mailto:hello@botb.com"
                      className="hover:text-botb-orange"
                    >
                      hello@botb.com
                    </a>
                  </dd>
                </div>
                <div>
                  <dt className="font-medium text-botb-text">Phone</dt>
                  <dd className="text-botb-muted">
                    <a href="tel:+442033184400" className="hover:text-botb-orange">
                      +44 (0)20 3318 4400
                    </a>
                  </dd>
                </div>
                <div>
                  <dt className="font-medium text-botb-text">Address</dt>
                  <dd className="text-botb-muted">
                    BOTB Ltd
                    <br />
                    London Gatwick
                    <br />
                    United Kingdom
                  </dd>
                </div>
                <div>
                  <dt className="font-medium text-botb-text">Opening hours</dt>
                  <dd className="text-botb-muted">
                    Monday to Friday, 9am – 5:30pm (GMT)
                  </dd>
                </div>
              </dl>
            </div>

            <div className="rounded-md border border-botb-card-border p-6">
              <h2 className="font-jost text-[20px] font-semibold text-botb-text">
                Frequently asked
              </h2>
              <ul className="mt-4 flex flex-col gap-2 text-[15px]">
                <li>
                  <Link
                    href="/how-to-play"
                    className="text-botb-muted hover:text-botb-orange"
                  >
                    How do I play and win?
                  </Link>
                </li>
                <li>
                  <Link
                    href="/how-to-play"
                    className="text-botb-muted hover:text-botb-orange"
                  >
                    How are winners chosen?
                  </Link>
                </li>
                <li>
                  <Link
                    href="/how-to-play"
                    className="text-botb-muted hover:text-botb-orange"
                  >
                    When are competitions drawn?
                  </Link>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
