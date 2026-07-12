import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function PricingPage() {
  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col overflow-x-hidden">
      <SiteNav
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Pricing" }]}
      />

      {/* HERO */}
      <main className="flex-1 flex flex-col items-center justify-center text-center max-w-3xl mx-auto w-full px-6 py-14 sm:py-20 lg:py-24">
        <span className="inline-block text-xs font-bold uppercase tracking-widest text-primary bg-primary/10 border border-primary/20 px-4 py-2 rounded-full mb-6">
          Kirmya Free Phase
        </span>

        <h1 className="font-extrabold tracking-tight leading-[1.05] text-4xl sm:text-5xl lg:text-6xl mb-5">
          Completely free for candidates &amp; recruiters.
        </h1>

        <p className="text-lg sm:text-xl leading-relaxed text-muted-foreground max-w-xl mb-10">
          During Kirmya&apos;s initial launch phase, the entire platform is
          completely free. All features — including AI profile builder,
          resume optimization, career coach chat, and talent sourcing tools —
          are accessible to everyone with zero fees.
        </p>

        <div className="bg-card border border-border/60 rounded-3xl p-9 max-w-lg w-full text-left shadow-sm">
          <h3 className="text-lg font-bold mb-3">
            What does this mean for you?
          </h3>
          <ul className="list-none space-y-2.5 text-sm text-muted-foreground leading-relaxed">
            <li>
              <strong className="text-foreground font-semibold">
                No Subscriptions:
              </strong>{" "}
              Free access to all Premium AI-features.
            </li>
            <li>
              <strong className="text-foreground font-semibold">
                No Hidden Fees:
              </strong>{" "}
              No payment gateway popups, locks, or limits.
            </li>
            <li>
              <strong className="text-foreground font-semibold">
                No Placement Fees:
              </strong>{" "}
              Recruiters can post jobs and hire without billing.
            </li>
            <li>
              <strong className="text-foreground font-semibold">
                Future-Ready Architecture:
              </strong>{" "}
              Designed cleanly to support modular billing in the future,
              without blocking current usage.
            </li>
          </ul>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
