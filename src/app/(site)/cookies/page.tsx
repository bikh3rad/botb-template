import { LegalPage } from "@/components/LegalPage";

const sections = [
  {
    heading: "What cookies are",
    body: [
      "Cookies are small text files placed on your device when you visit a website. They allow the site to remember your actions and preferences over a period of time, so you do not have to re-enter them each time you return.",
      "Similar technologies such as pixels and local storage perform comparable functions. In this policy we refer to all of these collectively as cookies.",
    ],
  },
  {
    heading: "Types of cookies we use",
    body: [
      "Essential cookies are required for the site to work, for example to keep you signed in and to process competition entries securely. These cannot be switched off without affecting core functionality.",
      "Performance and analytics cookies help us understand how visitors use the site so we can improve it. Functional cookies remember your choices, and where permitted, marketing cookies help us show relevant offers.",
    ],
  },
  {
    heading: "Managing cookies",
    body: [
      "You can manage non-essential cookies at any time through our cookie banner or your account settings. Most browsers also let you block or delete cookies through their own settings menus.",
      "Please note that disabling certain cookies may affect how the site functions and could prevent you from using some features.",
    ],
  },
  {
    heading: "Consent",
    body: [
      "When you first visit the site we ask for your consent to use non-essential cookies. Essential cookies are set automatically because they are necessary for the service to operate.",
      "Your consent choices are stored so we can respect them on future visits. You can review or change your preferences whenever you like through the cookie settings.",
    ],
  },
];

export default function CookiesPage() {
  return (
    <LegalPage title="Cookie Policy" updated="June 2026" sections={sections} />
  );
}
