import { LegalPage } from "@/components/LegalPage";

const sections = [
  {
    heading: "Eligibility",
    body: [
      "To enter any competition you must be at least 18 years of age and a resident of an eligible territory. By entering you confirm that you meet these requirements and that participation does not breach any law that applies to you.",
      "We reserve the right to request proof of age and identity before issuing any prize. Entries from individuals who do not satisfy the eligibility criteria may be declared void without refund.",
    ],
  },
  {
    heading: "How competitions work",
    body: [
      "Each competition has its own closing date, entry limit and prize, all of which are shown on the relevant competition page. Once a competition closes no further entries can be accepted.",
      "We may extend, shorten or withdraw a competition where circumstances beyond our reasonable control make this necessary. Where a competition is withdrawn before any winner is selected, entry fees paid for that competition will be refunded.",
    ],
  },
  {
    heading: "Skill element and Spot the Ball",
    body: [
      "Our Dream Car competitions include a genuine element of skill through the Spot the Ball game. Entrants must judge, using the players' body position and eye line, where the centre of the ball is most likely to be.",
      "An independent judge determines the winning coordinates after the competition closes. The judge's decision on the winning position is final and forms the basis on which the winner is selected.",
    ],
  },
  {
    heading: "Entries and tickets",
    body: [
      "Tickets are allocated once payment has been successfully processed. You may purchase as many tickets as you wish up to the maximum stated for the relevant competition.",
      "It is your responsibility to ensure the details associated with your entry are accurate. We are not liable for entries that are lost, delayed or incorrectly submitted due to factors outside our control.",
    ],
  },
  {
    heading: "Prizes and cash alternatives",
    body: [
      "The prize for each competition is described on its competition page. Where a cash alternative is offered, the winner may choose to receive the stated cash amount instead of the physical prize.",
      "Prizes are non-transferable and cannot be exchanged except as expressly provided. Any taxes, insurance, licensing or running costs associated with a prize are the responsibility of the winner unless stated otherwise.",
    ],
  },
  {
    heading: "Winner selection and liability",
    body: [
      "Winners are selected in accordance with the rules of each competition and are contacted directly using the details provided at entry. If we are unable to reach a winner within a reasonable period we may select an alternative winner.",
      "To the fullest extent permitted by law, our liability in connection with any competition is limited to the value of the entry fees paid. Nothing in these terms excludes liability that cannot lawfully be excluded.",
    ],
  },
];

export default function TermsPage() {
  return (
    <LegalPage title="Terms of Play" updated="June 2026" sections={sections} />
  );
}
