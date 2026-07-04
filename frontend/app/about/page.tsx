import Link from "next/link";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function AboutPage() {
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
        breadcrumb={[{ label: "Home", href: "/" }, { label: "About" }]}
      />

      {/* HERO */}
      <section
        style={{
          maxWidth: "840px",
          margin: "0 auto",
          padding: "clamp(56px,7vw,96px) 40px clamp(40px,5vw,56px)",
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
          Our story
        </div>
        <h1
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(38px,6vw,68px)",
            lineHeight: 1.02,
            letterSpacing: "-0.025em",
            margin: "0 0 24px",
          }}
        >
          A layoff is a sentence, not the&nbsp;story.
        </h1>
        <p
          style={{
            fontSize: "clamp(17px,2vw,20px)",
            lineHeight: 1.65,
            color: "#5B554C",
            maxWidth: "620px",
            margin: "0 auto",
          }}
        >
          Kirmya started the week our founder watched a 26-year industry veteran
          get filtered out by a résumé bot. We thought experience deserved
          better. So we built a place where it does.
        </p>
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
            background: "#2B2620",
            borderRadius: "24px",
            padding: "clamp(36px,4vw,52px)",
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit,minmax(180px,1fr))",
            gap: "32px",
            textAlign: "center",
          }}
        >
          <div>
            <div
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(34px,4.5vw,48px)",
                color: "#E7A57E",
              }}
            >
              2024
            </div>
            <div
              style={{ fontSize: "14px", color: "#C9C2B8", marginTop: "6px" }}
            >
              Founded
            </div>
          </div>
          <div>
            <div
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(34px,4.5vw,48px)",
                color: "#fff",
              }}
            >
              12,400+
            </div>
            <div
              style={{ fontSize: "14px", color: "#C9C2B8", marginTop: "6px" }}
            >
              Professionals placed
            </div>
          </div>
          <div>
            <div
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(34px,4.5vw,48px)",
                color: "#fff",
              }}
            >
              480
            </div>
            <div
              style={{ fontSize: "14px", color: "#C9C2B8", marginTop: "6px" }}
            >
              Hiring partners
            </div>
          </div>
          <div>
            <div
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(34px,4.5vw,48px)",
                color: "#A9C7B8",
              }}
            >
              100%
            </div>
            <div
              style={{ fontSize: "14px", color: "#C9C2B8", marginTop: "6px" }}
            >
              Reply guarantee
            </div>
          </div>
        </div>
      </section>

      {/* WHY WE EXIST */}
      <section
        style={{
          maxWidth: "1080px",
          margin: "0 auto",
          padding: "0 40px clamp(48px,6vw,72px)",
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit,minmax(300px,1fr))",
            gap: "40px",
            alignItems: "center",
          }}
        >
          <div>
            <div
              style={{
                fontSize: "13px",
                fontWeight: 700,
                letterSpacing: "0.12em",
                textTransform: "uppercase",
                color: "#C2683C",
                marginBottom: "14px",
              }}
            >
              Why we exist
            </div>
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(28px,3.6vw,40px)",
                lineHeight: 1.08,
                letterSpacing: "-0.02em",
                margin: "0 0 18px",
              }}
            >
              Hiring forgot that people are more than keywords.
            </h2>
            <p
              style={{
                fontSize: "16px",
                lineHeight: 1.7,
                color: "#4A443B",
                margin: "0 0 14px",
              }}
            >
              Recessions hit experienced workers hardest — and modern hiring
              tools make it worse, screening out exactly the judgment a company
              needs most when things are hard.
            </p>
            <p
              style={{ fontSize: "16px", lineHeight: 1.7, color: "#4A443B", margin: 0 }}
            >
              We&apos;re rebuilding the match around proof: what you&apos;ve
              shipped, the calls you&apos;ve made, the teams you&apos;ve
              steadied. Then we put a human on both sides of it.
            </p>
          </div>
          <div
            style={{
              background: "#F3ECE2",
              borderRadius: "22px",
              padding: "clamp(28px,3.5vw,40px)",
              display: "flex",
              flexDirection: "column",
              gap: "20px",
            }}
          >
            <div style={{ display: "flex", gap: "16px" }}>
              <span
                style={{
                  flex: "none",
                  width: "40px",
                  height: "40px",
                  borderRadius: "50%",
                  background: "rgba(79,124,106,0.16)",
                  color: "#4F7C6A",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "18px",
                }}
              >
                ◍
              </span>
              <div>
                <div
                  style={{
                    fontWeight: 600,
                    fontSize: "17px",
                    marginBottom: "3px",
                  }}
                >
                  Proof over pedigree
                </div>
                <div
                  style={{ fontSize: "15px", color: "#6B6357", lineHeight: 1.5 }}
                >
                  Outcomes lead. Logos and buzzwords don&apos;t.
                </div>
              </div>
            </div>
            <div style={{ display: "flex", gap: "16px" }}>
              <span
                style={{
                  flex: "none",
                  width: "40px",
                  height: "40px",
                  borderRadius: "50%",
                  background: "rgba(194,104,60,0.16)",
                  color: "#C2683C",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "18px",
                }}
              >
                ♥
              </span>
              <div>
                <div
                  style={{
                    fontWeight: 600,
                    fontSize: "17px",
                    marginBottom: "3px",
                  }}
                >
                  Dignity, always
                </div>
                <div
                  style={{ fontSize: "15px", color: "#6B6357", lineHeight: 1.5 }}
                >
                  No ghosting. Every applicant hears back.
                </div>
              </div>
            </div>
            <div style={{ display: "flex", gap: "16px" }}>
              <span
                style={{
                  flex: "none",
                  width: "40px",
                  height: "40px",
                  borderRadius: "50%",
                  background: "rgba(43,38,32,0.1)",
                  color: "#2B2620",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "18px",
                }}
              >
                ✦
              </span>
              <div>
                <div
                  style={{
                    fontWeight: 600,
                    fontSize: "17px",
                    marginBottom: "3px",
                  }}
                >
                  Free for candidates
                </div>
                <div
                  style={{ fontSize: "15px", color: "#6B6357", lineHeight: 1.5 }}
                >
                  Recruiters fund the platform. Job-seekers never pay.
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* PRINCIPLES */}
      <section
        style={{
          background: "#F3ECE2",
          borderTop: "1px solid #EFE7DC",
          borderBottom: "1px solid #EFE7DC",
        }}
      >
        <div
          style={{
            maxWidth: "1240px",
            margin: "0 auto",
            padding: "clamp(56px,6vw,84px) 40px",
          }}
        >
          <div
            style={{
              textAlign: "center",
              maxWidth: "560px",
              margin: "0 auto clamp(40px,5vw,52px)",
            }}
          >
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(30px,4vw,44px)",
                letterSpacing: "-0.02em",
                margin: 0,
              }}
            >
              What we stand on
            </h2>
          </div>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit,minmax(260px,1fr))",
              gap: "20px",
            }}
          >
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "18px",
                padding: "30px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "28px",
                  color: "#C2683C",
                  marginBottom: "14px",
                }}
              >
                01
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "19px",
                  margin: "0 0 8px",
                }}
              >
                Treat the gap with respect
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                Being between jobs is hard enough. The product should never make
                it feel worse.
              </p>
            </div>
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "18px",
                padding: "30px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "28px",
                  color: "#4F7C6A",
                  marginBottom: "14px",
                }}
              >
                02
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "19px",
                  margin: "0 0 8px",
                }}
              >
                Make experience legible
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                Help recruiters see the depth that a one-page résumé flattens.
              </p>
            </div>
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "18px",
                padding: "30px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "28px",
                  color: "#2B2620",
                  marginBottom: "14px",
                }}
              >
                03
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "19px",
                  margin: "0 0 8px",
                }}
              >
                Close the loop
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                A reply is the minimum. We guarantee one on every application.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* TEAM */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "clamp(56px,6vw,84px) 40px",
        }}
      >
        <div
          style={{
            textAlign: "center",
            maxWidth: "560px",
            margin: "0 auto clamp(36px,4vw,48px)",
          }}
        >
          <div
            style={{
              fontSize: "13px",
              fontWeight: 700,
              letterSpacing: "0.12em",
              textTransform: "uppercase",
              color: "#C2683C",
              marginBottom: "14px",
            }}
          >
            The team
          </div>
          <h2
            style={{
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 800,
              fontSize: "clamp(28px,3.6vw,40px)",
              letterSpacing: "-0.02em",
              margin: 0,
            }}
          >
            People who&apos;ve been there
          </h2>
        </div>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit,minmax(220px,1fr))",
            gap: "22px",
          }}
        >
          <div style={{ textAlign: "center" }}>
            <div
              style={{
                width: "96px",
                height: "96px",
                borderRadius: "50%",
                background: "#C2683C",
                margin: "0 auto 16px",
              }}
            />
            <div style={{ fontWeight: 600, fontSize: "17px" }}>Dana Reyes</div>
            <div style={{ fontSize: "14px", color: "#8A8175" }}>
              Founder &amp; CEO
            </div>
          </div>
          <div style={{ textAlign: "center" }}>
            <div
              style={{
                width: "96px",
                height: "96px",
                borderRadius: "50%",
                background: "#4F7C6A",
                margin: "0 auto 16px",
              }}
            />
            <div style={{ fontWeight: 600, fontSize: "17px" }}>Sam Okafor</div>
            <div style={{ fontSize: "14px", color: "#8A8175" }}>
              Head of Talent
            </div>
          </div>
          <div style={{ textAlign: "center" }}>
            <div
              style={{
                width: "96px",
                height: "96px",
                borderRadius: "50%",
                background: "#2B2620",
                margin: "0 auto 16px",
              }}
            />
            <div style={{ fontWeight: 600, fontSize: "17px" }}>Priya Nair</div>
            <div style={{ fontSize: "14px", color: "#8A8175" }}>
              Head of Product
            </div>
          </div>
          <div style={{ textAlign: "center" }}>
            <div
              style={{
                width: "96px",
                height: "96px",
                borderRadius: "50%",
                background: "#B08A2E",
                margin: "0 auto 16px",
              }}
            />
            <div style={{ fontWeight: 600, fontSize: "17px" }}>Marco Liu</div>
            <div style={{ fontSize: "14px", color: "#8A8175" }}>
              Head of Partnerships
            </div>
          </div>
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
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 800,
              color: "#fff",
              fontSize: "clamp(28px,4vw,44px)",
              lineHeight: 1.05,
              letterSpacing: "-0.02em",
              margin: "0 0 14px",
            }}
          >
            Join the comeback.
          </h2>
          <p
            style={{
              fontSize: "clamp(16px,2vw,18px)",
              color: "rgba(255,255,255,0.88)",
              margin: "0 auto 28px",
              maxWidth: "520px",
            }}
          >
            Whether you&apos;re hiring or looking, there&apos;s a more human
            way to do this.
          </p>
          <div
            style={{
              display: "flex",
              gap: "14px",
              justifyContent: "center",
              flexWrap: "wrap",
            }}
          >
            <Link
              href="/jobs"
              style={{
                background: "#fff",
                color: "#2B2620",
                fontSize: "16px",
                fontWeight: 600,
                padding: "16px 32px",
                borderRadius: "100px",
                textDecoration: "none",
                display: "inline-block",
              }}
            >
              Find work
            </Link>
            <Link
              href="/recruiter"
              style={{
                background: "rgba(255,255,255,0.14)",
                color: "#fff",
                fontSize: "16px",
                fontWeight: 600,
                padding: "16px 32px",
                borderRadius: "100px",
                border: "1px solid rgba(255,255,255,0.35)",
                textDecoration: "none",
                display: "inline-block",
              }}
            >
              Post a role
            </Link>
          </div>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
