import Link from "next/link";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function HomePage() {
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
      <SiteNav />

      {/* HERO */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "clamp(56px,8vw,104px) 40px clamp(44px,5vw,64px)",
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
            marginBottom: "28px",
          }}
        >
          Built for the moment between jobs
        </div>
        <h1
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(40px,7vw,76px)",
            lineHeight: 1.02,
            letterSpacing: "-0.025em",
            margin: "0 auto 24px",
            maxWidth: "900px",
          }}
        >
          You didn&apos;t lose your career.
          <br />
          <span style={{ color: "#C2683C" }}>You just lost that one job.</span>
        </h1>
        <p
          style={{
            fontSize: "clamp(17px,2vw,20px)",
            lineHeight: 1.6,
            color: "#5B554C",
            maxWidth: "640px",
            margin: "0 auto 36px",
          }}
        >
          The gap on your resume isn&apos;t a red flag — it&apos;s a chapter.
          Kirmya is where professionals like you come to breathe, regroup, and
          come back stronger. Real referrals. Real mentors. An AI coach that
          gets what you&apos;ve been through.
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
            href="/sign-in"
            style={{
              background: "#C2683C",
              color: "#fff",
              fontSize: "16px",
              fontWeight: 600,
              padding: "16px 32px",
              borderRadius: "100px",
              textDecoration: "none",
              display: "inline-block",
            }}
          >
            Start your comeback
          </Link>
          <Link
            href="/#how"
            style={{
              background: "#fff",
              color: "#2B2620",
              border: "1px solid #E2D8CB",
              fontSize: "16px",
              fontWeight: 600,
              padding: "16px 32px",
              borderRadius: "100px",
              textDecoration: "none",
              display: "inline-block",
            }}
          >
            See how it works
          </Link>
        </div>
        <div
          style={{
            marginTop: "24px",
            fontSize: "14px",
            color: "#8A8175",
          }}
        >
          Free to join · No spam · Your data stays yours
        </div>
      </section>

      {/* FEATURES */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "0 40px clamp(48px,6vw,80px)",
        }}
      >
        <div
          style={{
            textAlign: "center",
            maxWidth: "640px",
            margin: "0 auto clamp(36px,4vw,52px)",
          }}
        >
          <h2
            style={{
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 800,
              fontSize: "clamp(30px,4.4vw,46px)",
              lineHeight: 1.05,
              letterSpacing: "-0.02em",
              margin: "0 0 16px",
            }}
          >
            Everything the job hunt actually needs.
          </h2>
          <p
            style={{
              fontSize: "clamp(16px,2vw,18px)",
              lineHeight: 1.6,
              color: "#5B554C",
              margin: 0,
            }}
          >
            Not a feed to perform on. Not another platform to game. Just the
            tools, people, and guidance that get you hired.
          </p>
        </div>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit,minmax(300px,1fr))",
            gap: "20px",
          }}
        >
          <Link
            href="/resume"
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "32px",
              display: "block",
              textDecoration: "none",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "12px",
                background: "rgba(194,104,60,0.14)",
                marginBottom: "18px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "22px",
                color: "#C2683C",
              }}
            >
              ◈
            </div>
            <h3
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                margin: "0 0 8px",
              }}
            >
              AI Resume Coach
            </h3>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.55,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Upload your resume. Walk away with an ATS score, missing keywords,
              and edits that make recruiters stop scrolling.
            </p>
          </Link>
          <Link
            href="/referrals"
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "32px",
              display: "block",
              textDecoration: "none",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "12px",
                background: "rgba(79,124,106,0.14)",
                marginBottom: "18px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "22px",
                color: "#4F7C6A",
              }}
            >
              ⇄
            </div>
            <h3
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                margin: "0 0 8px",
              }}
            >
              Real Referrals
            </h3>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.55,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Skip the cold-apply black hole. Request warm intros from people
              already inside the companies you want.
            </p>
          </Link>
          <Link
            href="/mentorship"
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "32px",
              display: "block",
              textDecoration: "none",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "12px",
                background: "rgba(43,38,32,0.1)",
                marginBottom: "18px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "22px",
                color: "#2B2620",
              }}
            >
              ◍
            </div>
            <h3
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                margin: "0 0 8px",
              }}
            >
              Mentorship
            </h3>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.55,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Talk to someone who&apos;s sat in the seat you are aiming for.
              Book sessions. Get honest advice. Move faster.
            </p>
          </Link>
          <Link
            href="/communities"
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "32px",
              display: "block",
              textDecoration: "none",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "12px",
                background: "rgba(106,95,160,0.14)",
                marginBottom: "18px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "22px",
                color: "#6A5FA0",
              }}
            >
              ❖
            </div>
            <h3
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                margin: "0 0 8px",
              }}
            >
              Quiet Communities
            </h3>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.55,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Industry circles for ops, tech, HR, logistics, and more — where
              people share leads, not selfies.
            </p>
          </Link>
          <Link
            href="/career-paths"
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "32px",
              display: "block",
              textDecoration: "none",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "12px",
                background: "rgba(79,124,106,0.14)",
                marginBottom: "18px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "22px",
                color: "#4F7C6A",
              }}
            >
              ↗
            </div>
            <h3
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                margin: "0 0 8px",
              }}
            >
              Career Paths
            </h3>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.55,
                color: "#6B6357",
                margin: 0,
              }}
            >
              See where you can realistically go next, what it pays, and the
              exact skills standing between you and the offer.
            </p>
          </Link>
          <Link
            href="/coach"
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "32px",
              display: "block",
              textDecoration: "none",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "12px",
                background: "rgba(194,104,60,0.14)",
                marginBottom: "18px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "22px",
                color: "#C2683C",
              }}
            >
              ✦
            </div>
            <h3
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "20px",
                margin: "0 0 8px",
              }}
            >
              AI Career Coach
            </h3>
            <p
              style={{
                fontSize: "15px",
                lineHeight: 1.55,
                color: "#6B6357",
                margin: 0,
              }}
            >
              Interview prep, salary negotiation, &ldquo;what do I do this
              week?&rdquo; — your coach is always on, never judges, never gets
              tired.
            </p>
          </Link>
        </div>
      </section>
      {/* HOW IT WORKS */}
      <section
        id="how"
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
              maxWidth: "660px",
              margin: "0 auto clamp(40px,5vw,56px)",
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
              How Kirmya works
            </div>
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(30px,4vw,46px)",
                lineHeight: 1.05,
                letterSpacing: "-0.02em",
                margin: 0,
              }}
            >
              From lost to landed — here&apos;s the path.
            </h2>
          </div>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit,minmax(250px,1fr))",
              gap: "20px",
            }}
          >
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "20px",
                padding: "32px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "15px",
                  letterSpacing: "0.08em",
                  color: "#4F7C6A",
                  marginBottom: "18px",
                }}
              >
                01
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "20px",
                  margin: "0 0 10px",
                }}
              >
                Tell us where you are
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                One short profile, your resume, and what happened. No judgment
                here.
              </p>
            </div>
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "20px",
                padding: "32px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "15px",
                  letterSpacing: "0.08em",
                  color: "#4F7C6A",
                  marginBottom: "18px",
                }}
              >
                02
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "20px",
                  margin: "0 0 10px",
                }}
              >
                Get your personal playbook
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                Your AI playbook: resume score, skill gaps, target roles, and
                the three things to do <em>this week</em>.
              </p>
            </div>
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "20px",
                padding: "32px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "15px",
                  letterSpacing: "0.08em",
                  color: "#4F7C6A",
                  marginBottom: "18px",
                }}
              >
                03
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "20px",
                  margin: "0 0 10px",
                }}
              >
                Activate your network
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                Request referrals. Message mentors. Join the communities where
                real conversations are happening.
              </p>
            </div>
            <div
              style={{
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "20px",
                padding: "32px",
              }}
            >
              <div
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "15px",
                  letterSpacing: "0.08em",
                  color: "#4F7C6A",
                  marginBottom: "18px",
                }}
              >
                04
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "20px",
                  margin: "0 0 10px",
                }}
              >
                Track everything, stress less
              </h3>
              <p
                style={{
                  fontSize: "15px",
                  lineHeight: 1.55,
                  color: "#6B6357",
                  margin: 0,
                }}
              >
                Applications, referrals, sessions, and callbacks — all in one
                calm, clear dashboard. No chaos.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* ONE ACCOUNT, EVERY ROLE */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "clamp(8px,2vw,16px) 40px clamp(48px,6vw,80px)",
        }}
      >
        <div
          style={{
            background: "#2B2620",
            borderRadius: "24px",
            padding: "clamp(40px,5vw,64px) clamp(32px,4vw,56px)",
            position: "relative",
            overflow: "hidden",
          }}
        >
          <div
            style={{
              position: "absolute",
              top: "-120px",
              right: "-100px",
              width: "380px",
              height: "380px",
              borderRadius: "50%",
              background:
                "radial-gradient(circle, rgba(79,124,106,0.32), transparent 70%)",
              pointerEvents: "none",
            }}
          />
          <div
            style={{
              position: "relative",
              textAlign: "center",
              maxWidth: "700px",
              margin: "0 auto clamp(32px,4vw,44px)",
            }}
          >
            <div
              style={{
                display: "inline-block",
                fontSize: "13px",
                fontWeight: 700,
                letterSpacing: "0.1em",
                textTransform: "uppercase",
                color: "#E7A57E",
                background: "rgba(231,165,126,0.14)",
                padding: "8px 16px",
                borderRadius: "100px",
                marginBottom: "22px",
              }}
            >
              One account, every role
            </div>
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                color: "#fff",
                fontSize: "clamp(28px,4.2vw,46px)",
                lineHeight: 1.05,
                letterSpacing: "-0.02em",
                margin: "0 0 16px",
              }}
            >
              You&apos;re not just one thing.
              <br />
              Neither is Kirmya.
            </h2>
            <p
              style={{
                fontSize: "clamp(16px,2vw,18px)",
                lineHeight: 1.6,
                color: "rgba(255,255,255,0.78)",
                margin: 0,
              }}
            >
              The same free login lets you look for work, hire for your team,
              and mentor someone behind you — switch roles anytime, no second
              account.
            </p>
          </div>
          <div
            style={{
              position: "relative",
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit,minmax(240px,1fr))",
              gap: "16px",
            }}
          >
            <div
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                borderRadius: "18px",
                padding: "28px",
              }}
            >
              <div
                style={{
                  width: "46px",
                  height: "46px",
                  borderRadius: "12px",
                  background: "rgba(79,124,106,0.22)",
                  color: "#9FC7B5",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "22px",
                  marginBottom: "16px",
                }}
              >
                ◍
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "19px",
                  color: "#fff",
                  margin: "0 0 8px",
                }}
              >
                Job Seeker
              </h3>
              <p
                style={{
                  fontSize: "14px",
                  lineHeight: 1.55,
                  color: "rgba(255,255,255,0.7)",
                  margin: 0,
                }}
              >
                Get matched, request referrals, and track every application in
                one calm place.
              </p>
            </div>
            <div
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                borderRadius: "18px",
                padding: "28px",
              }}
            >
              <div
                style={{
                  width: "46px",
                  height: "46px",
                  borderRadius: "12px",
                  background: "rgba(194,104,60,0.22)",
                  color: "#E7A57E",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "22px",
                  marginBottom: "16px",
                }}
              >
                ✦
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "19px",
                  color: "#fff",
                  margin: "0 0 8px",
                }}
              >
                Recruiter
              </h3>
              <p
                style={{
                  fontSize: "14px",
                  lineHeight: 1.55,
                  color: "rgba(255,255,255,0.7)",
                  margin: 0,
                }}
              >
                Post roles and search reference-checked talent who are ready to
                start now.
              </p>
            </div>
            <div
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                borderRadius: "18px",
                padding: "28px",
              }}
            >
              <div
                style={{
                  width: "46px",
                  height: "46px",
                  borderRadius: "12px",
                  background: "rgba(255,255,255,0.12)",
                  color: "#fff",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "22px",
                  marginBottom: "16px",
                }}
              >
                ❋
              </div>
              <h3
                style={{
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 700,
                  fontSize: "19px",
                  color: "#fff",
                  margin: "0 0 8px",
                }}
              >
                Mentor
              </h3>
              <p
                style={{
                  fontSize: "14px",
                  lineHeight: 1.55,
                  color: "rgba(255,255,255,0.7)",
                  margin: 0,
                }}
              >
                Guide someone through the comeback you&apos;ve already made.
                Share what you know.
              </p>
            </div>
          </div>
          <div
            style={{
              position: "relative",
              textAlign: "center",
              marginTop: "clamp(28px,4vw,40px)",
            }}
          >
            <Link
                        href="/sign-in"
                        style={{
                          display: "inline-block",
                          background: "#fff",
                          color: "#2B2620",
                          fontSize: "16px",
                          fontWeight: 600,
                          padding: "17px 34px",
                          borderRadius: "100px",
                          textDecoration: "none",
                        }}
                      >
                        Create your free account →
                      </Link>
            <div
              style={{
                marginTop: "14px",
                fontSize: "14px",
                color: "rgba(255,255,255,0.55)",
              }}
            >
              Pick your roles during sign-up · change them anytime in Settings
            </div>
          </div>
        </div>
      </section>

      {/* TESTIMONIAL */}
      <section
        style={{
          maxWidth: "900px",
          margin: "0 auto",
          padding: "clamp(56px,6vw,84px) 40px",
          textAlign: "center",
        }}
      >
        <div
          style={{
            fontSize: "46px",
            lineHeight: 1,
            color: "#C2683C",
            fontFamily: "'Bricolage Grotesque',sans-serif",
            marginBottom: "14px",
          }}
        >
          &ldquo;
        </div>
        <p
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontWeight: 500,
            fontSize: "clamp(22px,3vw,32px)",
            lineHeight: 1.3,
            letterSpacing: "-0.01em",
            margin: "0 0 26px",
          }}
        >
          Twenty-two years in supply chain, gone in one email. Kirmya had me in
          front of three hiring managers in a week — and they were excited about
          exactly the experience I thought made me too expensive.
        </p>
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            gap: "14px",
          }}
        >
          <span
            style={{
              width: "48px",
              height: "48px",
              borderRadius: "50%",
              background: "#C2683C",
              display: "inline-block",
            }}
          />
          <div style={{ textAlign: "left" }}>
            <div style={{ fontWeight: 600, fontSize: "16px" }}>Marcus Hale</div>
            <div style={{ fontSize: "14px", color: "#8A8175" }}>
              Operations Director · placed in 11 days
            </div>
          </div>
        </div>
      </section>

      {/* CLOSING CTA */}
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
            padding: "clamp(44px,6vw,72px) clamp(40px,5vw,64px)",
            textAlign: "center",
          }}
        >
          <h2
            style={{
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 800,
              color: "#fff",
              fontSize: "clamp(28px,4.4vw,46px)",
              lineHeight: 1.05,
              letterSpacing: "-0.02em",
              margin: "0 auto 18px",
              maxWidth: "760px",
            }}
          >
            The hardest part is starting. We made that easy.
          </h2>
          <p
            style={{
              fontSize: "clamp(16px,2vw,18px)",
              lineHeight: 1.6,
              color: "rgba(255,255,255,0.9)",
              margin: "0 auto 32px",
              maxWidth: "620px",
            }}
          >
            Thousands of professionals have used Kirmya to move from uncertainty
            to an offer letter. Some came back stronger than before. The job you
            want is still out there — and your next step is just one quiet
            click.
          </p>
          <a
            href="/sign-in"
            style={{
              display: "inline-block",
              background: "#fff",
              color: "#2B2620",
              fontSize: "16px",
              fontWeight: 600,
              padding: "17px 34px",
              borderRadius: "100px",
              textDecoration: "none",
            }}
          >
            Create your free account →
          </a>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
