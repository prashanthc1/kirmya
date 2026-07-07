import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function PricingPage() {
  return (
    <div
      style={{
        background: "#FBF7F2",
        fontFamily: "'Public Sans',sans-serif",
        color: "#2B2620",
        minHeight: "100vh",
        overflowX: "hidden",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <SiteNav
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Pricing" }]}
      />

      {/* HERO */}
      <section
        style={{
          maxWidth: "840px",
          margin: "0 auto",
          padding: "clamp(56px,7vw,96px) 40px clamp(32px,4vw,48px)",
          textAlign: "center",
          flex: 1,
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <div
          style={{
            display: "inline-block",
            fontSize: "13px",
            fontWeight: 700,
            letterSpacing: "0.1em",
            textTransform: "uppercase",
            color: "#4F7C6A",
            background: "rgba(79,124,106,0.12)",
            padding: "8px 16px",
            borderRadius: "100px",
            marginBottom: "26px",
          }}
        >
          Kirmya Free Phase
        </div>
        <h1
          style={{
            fontFamily: "'Public Sans',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(38px,5.5vw,60px)",
            lineHeight: 1.05,
            letterSpacing: "-0.025em",
            margin: "0 0 20px",
            color: "#2B2620",
          }}
        >
          Completely free for candidates &amp; recruiters.
        </h1>
        <p
          style={{
            fontSize: "clamp(17px,2vw,20px)",
            lineHeight: 1.6,
            color: "#5B554C",
            maxWidth: "640px",
            margin: "0 0 40px",
          }}
        >
          During Kirmya's initial launch phase, the entire platform is completely free. All features — including AI profile builder, resume optimization, career coach chat, and talent sourcing tools — are accessible to everyone with zero fees.
        </p>

        <div
          style={{
            background: "#fff",
            border: "1px solid #EFE7DC",
            borderRadius: "24px",
            padding: "36px",
            maxWidth: "520px",
            width: "100%",
            textAlign: "left",
            boxShadow: "0 4px 20px rgba(43, 38, 32, 0.03)",
          }}
        >
          <h3 style={{ margin: "0 0 12px 0", fontSize: "18px", fontWeight: 700 }}>
            What does this mean for you?
          </h3>
          <ul style={{ paddingLeft: "20px", margin: 0, color: "#5B554C", lineHeight: 1.6, display: "flex", flexDirection: "column", gap: "10px" }}>
            <li><strong>No Subscriptions:</strong> Free access to all Premium AI-features.</li>
            <li><strong>No Hidden Fees:</strong> No payment gateway popups, locks, or limits.</li>
            <li><strong>No Placement Fees:</strong> Recruiters can post jobs and hire without billing.</li>
            <li><strong>Future-Ready Architecture:</strong> Designed cleanly to support modular billing in the future, without blocking current usage.</li>
          </ul>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
