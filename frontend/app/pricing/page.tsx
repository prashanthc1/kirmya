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
          Pricing for recruiters
        </div>
        <h1
          style={{
            fontFamily: "'Public Sans',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(38px,5.5vw,64px)",
            lineHeight: 1.02,
            letterSpacing: "-0.025em",
            margin: "0 0 18px",
          }}
        >
          Pay to hire, never to apply.
        </h1>
        <p
          style={{
            fontSize: "clamp(17px,2vw,20px)",
            lineHeight: 1.6,
            color: "#5B554C",
            maxWidth: "560px",
            margin: "0 auto",
          }}
        >
          Candidates use Kirmya free, forever. Recruiters fund the platform —
          with simple plans and no placement fees.
        </p>
      </section>

      {/* PLAN CARDS */}
      <section
        style={{
          maxWidth: "1180px",
          margin: "0 auto",
          padding: "0 40px clamp(40px,5vw,56px)",
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit,minmax(290px,1fr))",
            gap: "20px",
            alignItems: "stretch",
          }}
        >
          {/* STARTER */}
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "22px",
              padding: "34px",
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                marginBottom: "6px",
              }}
            >
              Starter
            </div>
            <p
              style={{
                fontSize: "14px",
                color: "#8A8175",
                margin: "0 0 22px",
                lineHeight: 1.5,
              }}
            >
              For one-off or occasional hires.
            </p>
            <div
              style={{
                display: "flex",
                alignItems: "baseline",
                gap: "6px",
                marginBottom: "4px",
              }}
            >
              <span
                style={{
                  fontFamily: "'Public Sans',sans-serif",
                  fontWeight: 800,
                  fontSize: "46px",
                  letterSpacing: "-0.02em",
                }}
              >
                $299
              </span>
              <span style={{ fontSize: "15px", color: "#8A8175" }}>/ role</span>
            </div>
            <div
              style={{
                fontSize: "13px",
                color: "#8A8175",
                marginBottom: "26px",
              }}
            >
              Billed per active posting
            </div>
            <a
              href="/sign-in"
              style={{
                display: "block",
                textAlign: "center",
                border: "1px solid #D8CFC2",
                color: "#2B2620",
                fontSize: "15px",
                fontWeight: 600,
                padding: "13px",
                borderRadius: "100px",
                marginBottom: "26px",
              }}
            >
              Start hiring
            </a>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "13px" }}
            >
              {[
                "1 active role posting",
                "Full candidate search",
                "Unlimited messages",
                "Reference-checked profiles",
              ].map((f) => (
                <div
                  key={f}
                  style={{ display: "flex", gap: "10px", fontSize: "15px", color: "#5B554C" }}
                >
                  <span style={{ color: "#4F7C6A", flex: "none" }}>✓</span> {f}
                </div>
              ))}
            </div>
          </div>

          {/* GROWTH */}
          <div
            style={{
              background: "#2B2620",
              color: "#fff",
              borderRadius: "22px",
              padding: "34px",
              display: "flex",
              flexDirection: "column",
              position: "relative",
              boxShadow: "0 18px 40px rgba(43,38,32,0.18)",
            }}
          >
            <span
              style={{
                position: "absolute",
                top: "22px",
                right: "22px",
                fontSize: "12px",
                fontWeight: 600,
                color: "#2B2620",
                background: "#E7A57E",
                padding: "5px 12px",
                borderRadius: "100px",
              }}
            >
              Most popular
            </span>
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                marginBottom: "6px",
              }}
            >
              Growth
            </div>
            <p
              style={{
                fontSize: "14px",
                color: "#C9C2B8",
                margin: "0 0 22px",
                lineHeight: 1.5,
              }}
            >
              For teams hiring throughout the year.
            </p>
            <div
              style={{
                display: "flex",
                alignItems: "baseline",
                gap: "6px",
                marginBottom: "4px",
              }}
            >
              <span
                style={{
                  fontFamily: "'Public Sans',sans-serif",
                  fontWeight: 800,
                  fontSize: "46px",
                  letterSpacing: "-0.02em",
                }}
              >
                $899
              </span>
              <span style={{ fontSize: "15px", color: "#C9C2B8" }}>
                / month
              </span>
            </div>
            <div
              style={{
                fontSize: "13px",
                color: "#C9C2B8",
                marginBottom: "26px",
              }}
            >
              Up to 5 active roles · billed annually
            </div>
            <a
              href="/sign-in"
              style={{
                display: "block",
                textAlign: "center",
                background: "#C2683C",
                color: "#fff",
                fontSize: "15px",
                fontWeight: 600,
                padding: "13px",
                borderRadius: "100px",
                marginBottom: "26px",
              }}
            >
              Start 14-day trial
            </a>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "13px" }}
            >
              {[
                "Up to 5 active roles",
                "Saved searches & shortlists",
                "Priority candidate matching",
                "3 team seats",
                "Analytics dashboard",
              ].map((f) => (
                <div
                  key={f}
                  style={{ display: "flex", gap: "10px", fontSize: "15px", color: "#E5DFD5" }}
                >
                  <span style={{ color: "#E7A57E", flex: "none" }}>✓</span> {f}
                </div>
              ))}
            </div>
          </div>

          {/* ENTERPRISE */}
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "22px",
              padding: "34px",
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                marginBottom: "6px",
              }}
            >
              Enterprise
            </div>
            <p
              style={{
                fontSize: "14px",
                color: "#8A8175",
                margin: "0 0 22px",
                lineHeight: 1.5,
              }}
            >
              For high-volume talent teams.
            </p>
            <div
              style={{
                display: "flex",
                alignItems: "baseline",
                gap: "6px",
                marginBottom: "4px",
              }}
            >
              <span
                style={{
                  fontFamily: "'Public Sans',sans-serif",
                  fontWeight: 800,
                  fontSize: "46px",
                  letterSpacing: "-0.02em",
                }}
              >
                Custom
              </span>
            </div>
            <div
              style={{
                fontSize: "13px",
                color: "#8A8175",
                marginBottom: "26px",
              }}
            >
              Volume pricing &amp; SSO
            </div>
            <a
              href="/about"
              style={{
                display: "block",
                textAlign: "center",
                border: "1px solid #D8CFC2",
                color: "#2B2620",
                fontSize: "15px",
                fontWeight: 600,
                padding: "13px",
                borderRadius: "100px",
                marginBottom: "26px",
              }}
            >
              Talk to sales
            </a>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "13px" }}
            >
              {[
                "Unlimited roles & seats",
                "Dedicated talent partner",
                "SSO & ATS integration",
                "Custom reporting",
              ].map((f) => (
                <div
                  key={f}
                  style={{ display: "flex", gap: "10px", fontSize: "15px", color: "#5B554C" }}
                >
                  <span style={{ color: "#4F7C6A", flex: "none" }}>✓</span> {f}
                </div>
              ))}
            </div>
          </div>
        </div>
        <div
          style={{
            textAlign: "center",
            marginTop: "24px",
            fontSize: "14px",
            color: "#8A8175",
          }}
        >
          No placement fees on any plan. Cancel anytime.
        </div>
      </section>

      {/* STATS BAR */}
      <section
        style={{
          maxWidth: "1080px",
          margin: "0 auto",
          padding: "0 40px clamp(48px,6vw,72px)",
        }}
      >
        <div
          style={{
            background: "#F3ECE2",
            border: "1px solid #EFE7DC",
            borderRadius: "22px",
            padding: "clamp(32px,4vw,48px)",
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit,minmax(220px,1fr))",
            gap: "32px",
            textAlign: "center",
          }}
        >
          <div>
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(32px,4vw,44px)",
                color: "#C2683C",
                letterSpacing: "-0.02em",
              }}
            >
              0%
            </div>
            <div
              style={{ fontSize: "15px", color: "#6B6357", marginTop: "6px" }}
            >
              Placement fees
            </div>
          </div>
          <div>
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(32px,4vw,44px)",
                color: "#4F7C6A",
                letterSpacing: "-0.02em",
              }}
            >
              8 wks
            </div>
            <div
              style={{ fontSize: "15px", color: "#6B6357", marginTop: "6px" }}
            >
              Average time to hire
            </div>
          </div>
          <div>
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(32px,4vw,44px)",
                color: "#2B2620",
                letterSpacing: "-0.02em",
              }}
            >
              100%
            </div>
            <div
              style={{ fontSize: "15px", color: "#6B6357", marginTop: "6px" }}
            >
              Reference-checked
            </div>
          </div>
        </div>
      </section>

      {/* PRICING FAQ */}
      <section
        style={{
          maxWidth: "760px",
          margin: "0 auto",
          padding: "0 40px clamp(48px,6vw,72px)",
        }}
      >
        <h2
          style={{
            fontFamily: "'Public Sans',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(26px,3.4vw,38px)",
            letterSpacing: "-0.02em",
            margin: "0 0 24px",
            textAlign: "center",
          }}
        >
          Pricing questions
        </h2>
        <div
          style={{ display: "flex", flexDirection: "column", gap: "14px" }}
        >
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "16px",
              padding: "24px",
            }}
          >
            <div
              style={{ fontWeight: 600, fontSize: "17px", marginBottom: "8px" }}
            >
              Are there really no placement fees?
            </div>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.6,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Correct. You pay for access to post and search — never a
              percentage of a hire&apos;s salary. The price is the price.
            </p>
          </div>
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "16px",
              padding: "24px",
            }}
          >
            <div
              style={{ fontWeight: 600, fontSize: "17px", marginBottom: "8px" }}
            >
              What counts as an &ldquo;active role&rdquo;?
            </div>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.6,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Any posting currently open to applicants. Close or fill a role
              and you can open another in its place at no extra cost.
            </p>
          </div>
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "16px",
              padding: "24px",
            }}
          >
            <div
              style={{ fontWeight: 600, fontSize: "17px", marginBottom: "8px" }}
            >
              Can I change plans later?
            </div>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.6,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Anytime, up or down. Upgrades apply immediately; downgrades take
              effect at your next billing cycle.
            </p>
          </div>
        </div>
        <div style={{ textAlign: "center", marginTop: "20px" }}>
          <a
            href="/faq"
            style={{
              color: "#C2683C",
              fontWeight: 600,
              fontSize: "15px",
            }}
          >
            See the full FAQ →
          </a>
        </div>
      </section>

      {/* CTA */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "0 40px clamp(56px,6vw,90px)",
        }}
      >
        <div
          style={{
            background: "#4F7C6A",
            borderRadius: "24px",
            padding: "clamp(40px,5vw,64px)",
            textAlign: "center",
          }}
        >
          <h2
            style={{
              fontFamily: "'Public Sans',sans-serif",
              fontWeight: 800,
              color: "#fff",
              fontSize: "clamp(28px,4vw,44px)",
              lineHeight: 1.05,
              letterSpacing: "-0.02em",
              margin: "0 0 14px",
            }}
          >
            Start hiring proven talent.
          </h2>
          <p
            style={{
              fontSize: "clamp(16px,2vw,18px)",
              color: "rgba(255,255,255,0.88)",
              margin: "0 auto 28px",
              maxWidth: "520px",
            }}
          >
            Try Growth free for 14 days. No card required to start.
          </p>
          <a
            href="/sign-in"
            style={{
              background: "#fff",
              color: "#2B2620",
              fontSize: "16px",
              fontWeight: 600,
              padding: "16px 34px",
              borderRadius: "100px",
              display: "inline-block",
            }}
          >
            Start free trial
          </a>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
