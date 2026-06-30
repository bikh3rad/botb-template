import { LegalPage } from "@/components/LegalPage";

const sections = [
  {
    heading: "How to make a complaint",
    body: [
      "We want every entrant to have a positive experience, but we know things can occasionally go wrong. If you are unhappy with any aspect of our service, please get in touch using the contact options on the site.",
      "To help us resolve your complaint quickly, please include your account details, a clear description of the issue and any reference numbers or screenshots that are relevant.",
    ],
  },
  {
    heading: "Our process",
    body: [
      "Once we receive your complaint, a member of our customer care team will review it and gather any further information needed. We treat every complaint seriously and aim to be fair, consistent and transparent throughout.",
      "We will keep a record of your complaint and our handling of it, and we will let you know who is dealing with your case.",
    ],
  },
  {
    heading: "Timescales",
    body: [
      "We aim to acknowledge your complaint within two working days of receiving it. Most issues are resolved within a short period, but more complex cases may take longer to investigate fully.",
      "Where a resolution is likely to take additional time, we will keep you updated on our progress and explain the reasons for any delay.",
    ],
  },
  {
    heading: "Escalation",
    body: [
      "If you are not satisfied with our initial response, you can ask for your complaint to be escalated to a senior member of the team for a further review.",
      "A fresh review will consider the original complaint, our response and any new information you provide, and we will share the outcome with you in writing.",
    ],
  },
  {
    heading: "Alternative dispute resolution",
    body: [
      "If we are unable to resolve your complaint to your satisfaction, you may be entitled to refer the matter to an independent alternative dispute resolution (ADR) provider.",
      "We will let you know which ADR scheme applies and how to contact them, so you can pursue your complaint through an impartial third party if you wish.",
    ],
  },
];

export default function ComplaintsPage() {
  return (
    <LegalPage
      title="Complaints Policy"
      updated="June 2026"
      sections={sections}
    />
  );
}
