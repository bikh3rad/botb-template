import { LegalPage } from "@/components/LegalPage";

const sections = [
  {
    heading: "Play responsibly",
    body: [
      "Entering competitions should always be fun and something you do for enjoyment, never as a way to make money or escape from worries. We want every entrant to stay in control and play within their means.",
      "If taking part ever stops feeling enjoyable, it may be a good moment to pause and take a break. There is no prize worth more than your wellbeing.",
    ],
  },
  {
    heading: "Setting limits",
    body: [
      "Before you play, it can help to decide how much time and money you are comfortable spending, and to treat that as a firm limit. Only ever spend what you can genuinely afford to lose.",
      "We offer tools that let you set personal limits on your activity. Setting a limit in advance makes it easier to keep your play balanced and stress-free.",
    ],
  },
  {
    heading: "Self-exclusion",
    body: [
      "If you feel you need a longer break, you can choose to self-exclude from your account for a set period. During this time you will not be able to enter competitions, and we will pause marketing communications.",
      "Self-exclusion is a positive step, and our support team can help you arrange it. You can ask us about the options that best suit your situation at any time.",
    ],
  },
  {
    heading: "Support resources",
    body: [
      "If you are worried about your own play or someone else's, free and confidential help is available from independent organisations that specialise in this area.",
      "We encourage anyone who feels their play may be becoming a problem to reach out for support early. Talking to someone can make a real difference, and help is always available.",
    ],
  },
  {
    heading: "Age verification",
    body: [
      "You must be at least 18 to take part, and we carry out checks to confirm the age of our entrants. Protecting young people is a responsibility we take seriously.",
      "If you share a device with anyone under 18, please keep your login details private and consider using parental controls to prevent underage access.",
    ],
  },
];

export default function MindfulPlayPage() {
  return (
    <LegalPage title="Mindful Play" updated="June 2026" sections={sections} />
  );
}
