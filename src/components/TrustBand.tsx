export function TrustBand() {
  return (
    <section className="bg-botb-gray py-12 text-center">
      <div className="mx-auto max-w-[1360px] px-4">
        <h2 className="font-jost text-[28px] font-normal text-[#9a9aa2] md:text-[40px]">
          EST. 1999 — £160M+ IN PRIZES
        </h2>
        <p className="mt-2 text-[18px] text-[#9a9aa2] md:text-[22px]">
          Guaranteed winners every week
        </p>
        <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
          <button className="rounded bg-botb-orange px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-white transition-colors hover:bg-botb-orange-hover">
            Dream Car »
          </button>
          <button className="rounded border border-botb-card-border bg-white px-6 py-2.5 font-jost text-[15px] font-medium uppercase text-botb-text transition-colors hover:border-botb-orange hover:text-botb-orange">
            First Visit? »
          </button>
        </div>
      </div>
    </section>
  );
}
