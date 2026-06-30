export function JustLaunched() {
  return (
    <section className="bg-white">
      <div className="mx-auto grid max-w-[1360px] grid-cols-1 items-center gap-6 px-4 py-10 md:grid-cols-2">
        <div className="overflow-hidden rounded-lg">
          <img
            src="/images/misc/just-launched.webp"
            alt="Just launched competitions"
            className="h-full w-full object-cover"
          />
        </div>
        <div className="text-center">
          <p className="text-[13px] font-medium uppercase tracking-wide text-botb-muted">
            Just Launched!
          </p>
          <h2 className="mt-2 font-jost text-[30px] font-bold text-botb-text md:text-[38px]">
            Win cars, bikes, tech or cash!
          </h2>
          <button className="mt-6 rounded bg-botb-orange px-8 py-2.5 font-jost text-[15px] font-medium uppercase text-white transition-colors hover:bg-botb-orange-hover">
            Enter Now
          </button>
        </div>
      </div>
    </section>
  );
}
