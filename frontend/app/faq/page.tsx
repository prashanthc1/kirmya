import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function FaqPage() {
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
        breadcrumb={[{ label: "Home", href: "/" }, { label: "FAQ" }]}
      />

      {/* HERO */}
      <section
        style={{
          maxWidth: "760px",
          margin: "0 auto",
          padding: "clamp(56px,7vw,96px) 40px clamp(36px,4vw,48px)",
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
          Help center
        </div>
        <h1
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(38px,5.5vw,64px)",
            lineHeight: 1.02,
            letterSpacing: "-0.025em",
            margin: "0 0 18px",
          }}
        >
          Questions, answered.
        </h1>
        <p
          style={{
            fontSize: "clamp(17px,2vw,20px)",
            lineHeight: 1.6,
            color: "#5B554C",
            margin: 0,
          }}
        >
          Everything you need to know about finding work and hiring on Kirmya.
        </p>
      </section>

      {/* FAQ BODY */}
      <section
        style={{
          maxWidth: "1080px",
          margin: "0 auto",
          padding: "0 40px clamp(56px,6vw,90px)",
          display: "grid",
          gridTemplateColumns: "220px 1fr",
          gap: "40px",
          alignItems: "start",
        }}
      >
        {/* SIDEBAR NAV */}
        <aside
          style={{
            position: "sticky",
            top: "96px",
            display: "flex",
            flexDirection: "column",
            gap: "6px",
          }}
        >
          <div
            style={{
              fontSize: "13px",
              fontWeight: 700,
              letterSpacing: "0.06em",
              textTransform: "uppercase",
              color: "#8A8175",
              marginBottom: "8px",
            }}
          >
            Topics
          </div>
          <a
            href="#candidates"
            style={{
              fontSize: "15px",
              fontWeight: 600,
              color: "#C2683C",
              padding: "8px 0",
            }}
          >
            For candidates
          </a>
          <a
            href="#recruiters"
            style={{ fontSize: "15px", color: "#5B554C", padding: "8px 0" }}
          >
            For recruiters
          </a>
          <a
            href="#account"
            style={{ fontSize: "15px", color: "#5B554C", padding: "8px 0" }}
          >
            Account &amp; privacy
          </a>
          <a
            href="#pricing"
            style={{ fontSize: "15px", color: "#5B554C", padding: "8px 0" }}
          >
            Pricing
          </a>
        </aside>

        {/* FAQ SECTIONS */}
        <div style={{ display: "flex", flexDirection: "column", gap: "40px" }}>
          <div id="candidates">
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "24px",
                letterSpacing: "-0.01em",
                margin: "0 0 16px",
              }}
            >
              For candidates
            </h2>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "12px" }}
            >
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    Is Kirmya really free for job-seekers?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  Yes — completely. Creating a profile, browsing roles, and
                  applying are always free. Recruiters and employers fund the
                  platform, so you never pay to find work.
                </p>
              </div>
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    I was laid off recently. Will that count against me?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  Never. Kirmya was built for exactly this moment. Recruiters
                  here understand that a layoff is about budgets and timing, not
                  your ability. Your profile leads with your work, not your
                  employment gaps.
                </p>
              </div>
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    Will I actually hear back when I apply?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  Always. Every application gets a real response within five
                  business days — it&apos;s a guarantee we hold recruiters to.
                  No silent rejections, no black holes.
                </p>
              </div>
            </div>
          </div>

          <div id="recruiters">
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "24px",
                letterSpacing: "-0.01em",
                margin: "0 0 16px",
              }}
            >
              For recruiters
            </h2>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "12px" }}
            >
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    How are candidates vetted?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  Every professional completes an outcomes-based profile and is
                  reference-checked before appearing in search. You see verified
                  track records, not self-reported buzzwords.
                </p>
              </div>
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    How fast can I expect to hire?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  87% of roles posted on Kirmya fill within eight weeks. Because
                  candidates are pre-vetted and ready to start, your shortlist
                  is shorter and stronger from day one.
                </p>
              </div>
            </div>
          </div>

          <div id="account">
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "24px",
                letterSpacing: "-0.01em",
                margin: "0 0 16px",
              }}
            >
              Account &amp; privacy
            </h2>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "12px" }}
            >
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    Can my current employer see my profile?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  You control your visibility. Private mode hides your profile
                  from specific companies and shows it only to recruiters you
                  choose to engage. Your data is never sold.
                </p>
              </div>
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    How do I delete my account?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  From Settings → Account, one click permanently removes your
                  profile and all associated data within 48 hours. No retention
                  games.
                </p>
              </div>
            </div>
          </div>

          <div id="pricing">
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "24px",
                letterSpacing: "-0.01em",
                margin: "0 0 16px",
              }}
            >
              Pricing
            </h2>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "12px" }}
            >
              <div
                style={{
                  background: "#fff",
                  border: "1px solid #EFE7DC",
                  borderRadius: "16px",
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    padding: "22px 24px 8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: "16px",
                  }}
                >
                  <span style={{ fontWeight: 600, fontSize: "17px", color: "#2B2620" }}>
                    What does it cost to post roles?
                  </span>
                  <span style={{ flex: "none", fontSize: "22px", color: "#C2683C" }}>–</span>
                </div>
                <p style={{ margin: 0, padding: "0 24px 22px", fontSize: "16px", lineHeight: 1.65, color: "#5B554C" }}>
                  Recruiters pay per active role or via an annual seat plan for
                  high-volume teams. There are no placement fees and no surprise
                  charges — talk to us and we&apos;ll size a plan to your
                  hiring.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* STILL HAVE A QUESTION */}
      <section
        style={{
          maxWidth: "1080px",
          margin: "0 auto",
          padding: "0 40px clamp(56px,6vw,90px)",
        }}
      >
        <div
          style={{
            background: "#F3ECE2",
            border: "1px solid #EFE7DC",
            borderRadius: "22px",
            padding: "clamp(36px,4vw,52px)",
            textAlign: "center",
          }}
        >
          <h2
            style={{
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 800,
              fontSize: "clamp(24px,3vw,32px)",
              letterSpacing: "-0.02em",
              margin: "0 0 10px",
            }}
          >
            Still have a question?
          </h2>
          <p style={{ fontSize: "16px", color: "#6B6357", margin: "0 0 24px" }}>
            Our team replies to every message — usually within a day.
          </p>
          <a
            href="/about"
            style={{
              background: "#C2683C",
              color: "#fff",
              fontSize: "16px",
              fontWeight: 600,
              padding: "15px 32px",
              borderRadius: "100px",
              display: "inline-block",
            }}
          >
            Contact support
          </a>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
