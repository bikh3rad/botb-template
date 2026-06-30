interface LegalSection {
  heading: string;
  body: string[];
}

interface LegalPageProps {
  title: string;
  updated?: string;
  sections: LegalSection[];
}

export function LegalPage({ title, updated, sections }: LegalPageProps) {
  return (
    <div className="mx-auto max-w-4xl px-4 py-10">
      <div className="mb-8 rounded-lg bg-botb-gray p-4 text-sm text-botb-muted">
        This is a demo clone of botb.com for portfolio purposes — the text below
        is placeholder and not legal advice.
      </div>

      <h1 className="font-jost text-[28px] font-bold uppercase text-botb-text md:text-[36px]">
        {title}
      </h1>
      {updated ? (
        <p className="mt-2 text-sm text-botb-muted">Last updated: {updated}</p>
      ) : null}

      {sections.map((section) => (
        <section key={section.heading}>
          <h2 className="mb-2 mt-8 font-jost text-[20px] font-semibold text-botb-text">
            {section.heading}
          </h2>
          <div className="space-y-3 text-[15px] leading-relaxed text-botb-muted">
            {section.body.map((paragraph, index) => (
              <p key={index}>{paragraph}</p>
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}
