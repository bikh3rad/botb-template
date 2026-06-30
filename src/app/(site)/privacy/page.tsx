import { LegalPage } from "@/components/LegalPage";

const sections = [
  {
    heading: "Data we collect",
    body: [
      "When you create an account or enter a competition we collect information such as your name, date of birth, email address, postal address and payment details. We also collect technical data including your IP address, device type and how you interact with our site.",
      "We only collect the information we need to operate the service, verify eligibility and keep your account secure. You can choose not to provide certain details, although this may limit your ability to enter competitions.",
    ],
  },
  {
    heading: "How we use your information",
    body: [
      "We use your data to process entries, allocate tickets, contact winners and provide customer support. We also use it to detect fraud, meet our legal obligations and improve the products and features we offer.",
      "Where you have given consent, we may send you marketing communications about competitions and offers. You can withdraw this consent at any time using the unsubscribe link in our emails or by updating your account preferences.",
    ],
  },
  {
    heading: "Cookies",
    body: [
      "We use cookies and similar technologies to keep you signed in, remember your preferences and understand how the site is used. Some cookies are essential for the site to function, while others help us measure performance.",
      "You can control non-essential cookies through your browser settings or our cookie banner. For full details please see our Cookie Policy.",
    ],
  },
  {
    heading: "Third parties",
    body: [
      "We share data with trusted service providers who help us run the platform, such as payment processors, hosting providers and analytics partners. These providers may only use your data to perform services on our behalf.",
      "We do not sell your personal data. We may disclose information where required by law, to enforce our terms, or to protect the rights and safety of our users and our business.",
    ],
  },
  {
    heading: "Your rights",
    body: [
      "Depending on where you live, you may have the right to access, correct, delete or restrict the processing of your personal data, and to object to certain uses. You may also have the right to request a copy of your data in a portable format.",
      "To exercise any of these rights, contact us using the details below. We will respond within the timeframe required by applicable law and may need to verify your identity first.",
    ],
  },
  {
    heading: "Contact",
    body: [
      "If you have questions about this policy or how we handle your data, you can reach our team through the contact options provided on the site. We aim to acknowledge privacy enquiries promptly.",
      "If you are not satisfied with our response, you may have the right to lodge a complaint with your local data protection authority.",
    ],
  },
];

export default function PrivacyPage() {
  return (
    <LegalPage title="Privacy Policy" updated="June 2026" sections={sections} />
  );
}
