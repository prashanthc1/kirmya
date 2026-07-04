import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function MentorsPage() {
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
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Mentors" }]}
      />

      {/* HEADING */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "clamp(40px,5vw,64px) 40px clamp(20px,3vw,28px)",
        }}
      >
        <h1
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(34px,5vw,56px)",
            lineHeight: 1.02,
            letterSpacing: "-0.025em",
            margin: "0 0 14px",
          }}
        >
          Mentors who&apos;ve walked the road
        </h1>
        <p
          style={{
            fontSize: "clamp(16px,2vw,19px)",
            lineHeight: 1.6,
            color: "#5B554C",
            maxWidth: "580px",
            margin: 0,
          }}
        >
          2,100+ experienced professionals offering guidance. Filter by focus
          or availability.
        </p>
      </section>

      {/* SEARCH */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "0 40px 24px",
        }}
      >
        <div
          style={{
            background: "#fff",
            border: "1px solid #EFE7DC",
            borderRadius: "16px",
            padding: "14px",
            display: "flex",
            gap: "12px",
            flexWrap: "wrap",
            alignItems: "center",
          }}
        >
          <div
            style={{
              flex: "1 1 280px",
              display: "flex",
              alignItems: "center",
              gap: "10px",
              padding: "0 14px",
            }}
          >
            <span style={{ color: "#8A8175", fontSize: "18px" }}>⌕</span>
            <input
              placeholder="Field, focus area, or name"
              style={{
                border: "none",
                outline: "none",
                background: "transparent",
                fontSize: "16px",
                color: "#2B2620",
                width: "100%",
                padding: "12px 0",
              }}
            />
          </div>
          <button
            style={{
              border: "none",
              background: "#C2683C",
              color: "#fff",
              fontSize: "15px",
              fontWeight: 600,
              padding: "13px 28px",
              borderRadius: "100px",
              cursor: "pointer",
            }}
          >
            Search
          </button>
        </div>
        <div
          style={{
            display: "flex",
            gap: "10px",
            flexWrap: "wrap",
            marginTop: "16px",
          }}
        >
          {["Operations", "Engineering", "Product", "Available now"].map(
            (tag) => (
              <button
                key={tag}
                style={{
                  border: "1px solid #E2D9CC",
                  background: "#fff",
                  color: "#5B554C",
                  fontSize: "13px",
                  fontWeight: 500,
                  padding: "7px 16px",
                  borderRadius: "100px",
                  cursor: "pointer",
                }}
              >
                {tag}
              </button>
            )
          )}
        </div>
      </section>

      {/* MENTOR CARDS */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "8px 40px clamp(56px,6vw,90px)",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            marginBottom: "20px",
            flexWrap: "wrap",
            gap: "12px",
          }}
        >
          <span style={{ fontSize: "15px", color: "#6B6357" }}>
            <strong style={{ color: "#2B2620" }}>2,148</strong> mentors
          </span>
          <button
            style={{
              border: "none",
              background: "none",
              cursor: "pointer",
              fontFamily: "'Public Sans',sans-serif",
              fontSize: "14px",
              color: "#5B554C",
            }}
          >
            Sort: <strong style={{ color: "#2B2620" }}>Top rated</strong> ▾
          </button>
        </div>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fill,minmax(300px,1fr))",
            gap: "18px",
          }}
        >
          {/* MENTOR 1 */}
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "26px",
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "14px",
                marginBottom: "16px",
              }}
            >
              <span
                style={{
                  width: "54px",
                  height: "54px",
                  borderRadius: "50%",
                  background: "#C2683C",
                  flex: "none",
                  display: "inline-block",
                }}
              />
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ fontWeight: 600, fontSize: "17px" }}>
                  Daniela Cruz
                </div>
                <div style={{ fontSize: "13px", color: "#8A8175" }}>
                  VP Operations · Maersk
                </div>
              </div>
            </div>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                marginBottom: "14px",
                fontSize: "13px",
              }}
            >
              <span style={{ color: "#C2683C", fontWeight: 700 }}>★ 4.9</span>
              <span style={{ color: "#8A8175" }}>· 214 sessions</span>
            </div>
            <div
              style={{
                display: "flex",
                gap: "7px",
                flexWrap: "wrap",
                marginBottom: "18px",
              }}
            >
              <span
                style={{
                  fontSize: "12px",
                  color: "#5B554C",
                  background: "#F3ECE2",
                  padding: "5px 12px",
                  borderRadius: "100px",
                }}
              >
                Supply chain
              </span>
              <span
                style={{
                  fontSize: "12px",
                  color: "#5B554C",
                  background: "#F3ECE2",
                  padding: "5px 12px",
                  borderRadius: "100px",
                }}
              >
                Leadership
              </span>
            </div>
            <div
              style={{
                marginTop: "auto",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: "10px",
              }}
            >
              <span
                style={{
                  fontSize: "13px",
                  color: "#4F7C6A",
                  fontWeight: 600,
                }}
              >
                Available this week
              </span>
              <button
                style={{
                  border: "none",
                  background: "none",
                  cursor: "pointer",
                  fontFamily: "'Public Sans',sans-serif",
                  fontSize: "14px",
                  color: "#C2683C",
                  fontWeight: 600,
                }}
              >
                Request →
              </button>
            </div>
          </div>

          {/* MENTOR 2 */}
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "26px",
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "14px",
                marginBottom: "16px",
              }}
            >
              <span
                style={{
                  width: "54px",
                  height: "54px",
                  borderRadius: "50%",
                  background: "#4F7C6A",
                  flex: "none",
                  display: "inline-block",
                }}
              />
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ fontWeight: 600, fontSize: "17px" }}>
                  Sam Okafor
                </div>
                <div style={{ fontSize: "13px", color: "#8A8175" }}>
                  Staff Engineer · Stripe
                </div>
              </div>
            </div>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                marginBottom: "14px",
                fontSize: "13px",
              }}
            >
              <span style={{ color: "#C2683C", fontWeight: 700 }}>★ 5.0</span>
              <span style={{ color: "#8A8175" }}>· 96 sessions</span>
            </div>
            <div
              style={{
                display: "flex",
                gap: "7px",
                flexWrap: "wrap",
                marginBottom: "18px",
              }}
            >
              <span
                style={{
                  fontSize: "12px",
                  color: "#5B554C",
                  background: "#F3ECE2",
                  padding: "5px 12px",
                  borderRadius: "100px",
                }}
              >
                Backend
              </span>
              <span
                style={{
                  fontSize: "12px",
                  color: "#5B554C",
                  background: "#F3ECE2",
                  padding: "5px 12px",
                  borderRadius: "100px",
                }}
              >
                System design
              </span>
            </div>
            <div
              style={{
                marginTop: "auto",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: "10px",
              }}
            >
              <span
                style={{
                  fontSize: "13px",
                  color: "#8A8175",
                  fontWeight: 600,
                }}
              >
                Next: in 5 days
              </span>
              <button
                style={{
                  border: "none",
                  background: "none",
                  cursor: "pointer",
                  fontFamily: "'Public Sans',sans-serif",
                  fontSize: "14px",
                  color: "#C2683C",
                  fontWeight: 600,
                }}
              >
                Request →
              </button>
            </div>
          </div>

          {/* MENTOR 3 */}
          <div
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              padding: "26px",
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "14px",
                marginBottom: "16px",
              }}
            >
              <span
                style={{
                  width: "54px",
                  height: "54px",
                  borderRadius: "50%",
                  background: "#6A5FA0",
                  flex: "none",
                  display: "inline-block",
                }}
              />
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ fontWeight: 600, fontSize: "17px" }}>
                  Mei-Ling Tan
                </div>
                <div style={{ fontSize: "13px", color: "#8A8175" }}>
                  Director of Product · Figma
                </div>
              </div>
            </div>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                marginBottom: "14px",
                fontSize: "13px",
              }}
            >
              <span style={{ color: "#C2683C", fontWeight: 700 }}>★ 4.8</span>
              <span style={{ color: "#8A8175" }}>· 180 sessions</span>
            </div>
            <div
              style={{
                display: "flex",
                gap: "7px",
                flexWrap: "wrap",
                marginBottom: "18px",
              }}
            >
              <span
                style={{
                  fontSize: "12px",
                  color: "#5B554C",
                  background: "#F3ECE2",
                  padding: "5px 12px",
                  borderRadius: "100px",
                }}
              >
                Product
              </span>
              <span
                style={{
                  fontSize: "12px",
                  color: "#5B554C",
                  background: "#F3ECE2",
                  padding: "5px 12px",
                  borderRadius: "100px",
                }}
              >
                Career pivots
              </span>
            </div>
            <div
              style={{
                marginTop: "auto",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: "10px",
              }}
            >
              <span
                style={{
                  fontSize: "13px",
                  color: "#4F7C6A",
                  fontWeight: 600,
                }}
              >
                Available this week
              </span>
              <button
                style={{
                  border: "none",
                  background: "none",
                  cursor: "pointer",
                  fontFamily: "'Public Sans',sans-serif",
                  fontSize: "14px",
                  color: "#C2683C",
                  fontWeight: 600,
                }}
              >
                Request →
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* BECOME A MENTOR CTA */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "0 40px clamp(56px,6vw,90px)",
        }}
      >
        <div
          style={{
            background: "#2B2620",
            color: "#fff",
            borderRadius: "24px",
            padding: "clamp(36px,4vw,52px)",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: "30px",
            flexWrap: "wrap",
          }}
        >
          <div style={{ maxWidth: "560px" }}>
            <h2
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 800,
                fontSize: "clamp(24px,3vw,34px)",
                lineHeight: 1.08,
                letterSpacing: "-0.02em",
                margin: "0 0 8px",
              }}
            >
              Been through it yourself?
            </h2>
            <p style={{ fontSize: "16px", color: "#C9C2B8", margin: 0 }}>
              Join 2,100+ mentors giving back on their own schedule.
            </p>
          </div>
          <a
            href="/sign-in"
            style={{
              background: "#fff",
              color: "#2B2620",
              fontSize: "16px",
              fontWeight: 600,
              padding: "15px 32px",
              borderRadius: "100px",
              whiteSpace: "nowrap",
            }}
          >
            Become a mentor
          </a>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
