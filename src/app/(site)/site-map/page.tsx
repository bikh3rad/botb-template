import Link from "next/link";

type LinkItem = {
  label: string;
  href: string;
};

type LinkGroup = {
  heading: string;
  links: LinkItem[];
};

const groups: LinkGroup[] = [
  {
    heading: "Competitions",
    links: [
      { label: "Home", href: "/" },
      { label: "Competitions", href: "/competitions" },
      { label: "Dream Car", href: "/prizes/cars" },
      { label: "Winners", href: "/winners" },
      { label: "Prize Collections", href: "/collections/prize-collections" },
      { label: "How to Play", href: "/how-to-play" },
      { label: "BOTB Pass", href: "/botb-pass" },
      { label: "Testimonials", href: "/testimonials" },
    ],
  },
  {
    heading: "Account",
    links: [
      { label: "Login", href: "/login" },
      { label: "Register", href: "/register" },
      { label: "Account", href: "/account" },
      { label: "Cart", href: "/cart" },
      { label: "Perks for Playing", href: "/account/perks-for-playing" },
      { label: "Affiliates", href: "/affiliates" },
    ],
  },
  {
    heading: "About",
    links: [
      { label: "History & Locations", href: "/about/history-locations" },
      { label: "Charity", href: "/about/charity" },
      { label: "Blog", href: "/under-the-hood" },
      { label: "Contact", href: "/contact-us" },
    ],
  },
  {
    heading: "Legal",
    links: [
      { label: "Terms", href: "/terms" },
      { label: "Privacy", href: "/privacy" },
      { label: "Cookies", href: "/cookies" },
      { label: "Complaints", href: "/complaints" },
      { label: "Mindful Play", href: "/mindful-play" },
    ],
  },
];

export default function SiteMapPage() {
  return (
    <section className="bg-white">
      <div className="mx-auto max-w-[1360px] px-4 py-10">
        <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
          Site Map
        </h1>

        <div className="mt-8 grid grid-cols-2 gap-8 md:grid-cols-4">
          {groups.map((group) => (
            <div key={group.heading}>
              <h2 className="font-jost text-[16px] font-semibold uppercase text-botb-text">
                {group.heading}
              </h2>
              <ul className="mt-4 flex flex-col gap-2">
                {group.links.map((link) => (
                  <li key={link.href + link.label}>
                    <Link
                      href={link.href}
                      className="text-[15px] text-botb-muted hover:text-botb-orange"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
