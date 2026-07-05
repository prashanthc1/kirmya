import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function ResumePage() {
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
        breadcrumb={[
          { label: "Home", href: "/" },
          { label: "Dashboard", href: "/dashboard" },
          { label: "Resume Check" },
        ]}
      />

      {/* HEADING */}
      <section
        style={{
          maxWidth: "1080px",
          margin: "0 auto",
          padding: "clamp(36px,4vw,52px) 40px clamp(8px,2vw,16px)",
        }}
      >
        <div
          style={{
            fontSize: "13px",
            fontWeight: 700,
            letterSpacing: "0.12em",
            textTransform: "uppercase",
            color: "#C2683C",
            marginBottom: "12px",
          }}
        >
          Resume tools
        </div>
        <h1
          style={{
            fontFamily: "'Public Sans',sans-serif",
            fontWeight: 800,
            fontSize: "clamp(32px,4.5vw,52px)",
            lineHeight: 1.02,
            letterSpacing: "-0.025em",
            margin: "0 0 18px",
          }}
        >
          Beat the bots. Reach a human.
        </h1>
        <p
          style={{
            fontSize: "clamp(16px,2vw,19px)",
            lineHeight: 1.6,
            color: "#5B554C",
            maxWidth: "600px",
            margin: "0 0 24px",
          }}
        >
          Most résumés are filtered by an ATS before anyone reads them. Check
          yours in seconds and see exactly what to fix.
        </p>
        <div
          style={{
            display: "inline-flex",
            gap: "6px",
            background: "#F3ECE2",
            borderRadius: "100px",
            padding: "5px",
          }}
        >
          <span
            style={{
              fontSize: "14px",
              fontWeight: 600,
              color: "#fff",
              background: "#C2683C",
              padding: "9px 20px",
              borderRadius: "100px",
            }}
          >
            ATS Checker
          </span>
          <a
            href="/resume/builder"
            style={{
              fontSize: "14px",
              fontWeight: 600,
              color: "#5B554C",
              padding: "9px 20px",
              borderRadius: "100px",
            }}
          >
            Resume Builder
          </a>
        </div>
      </section>

      {/* UPLOAD + CHECKLIST */}
      <section
        style={{
          maxWidth: "1080px",
          margin: "0 auto",
          padding: "clamp(20px,3vw,28px) 40px clamp(56px,6vw,90px)",
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "1.3fr 0.7fr",
            gap: "24px",
            alignItems: "start",
          }}
        >
          {/* DROPZONE */}
          <div
            style={{
              background: "#fff",
              border: "2px dashed #D8CFC2",
              borderRadius: "20px",
              padding: "clamp(40px,6vw,72px) 40px",
              textAlign: "center",
              cursor: "pointer",
            }}
          >
            <div
              style={{
                width: "70px",
                height: "70px",
                borderRadius: "18px",
                background: "rgba(194,104,60,0.12)",
                color: "#C2683C",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "30px",
                margin: "0 auto 22px",
              }}
            >
              ↑
            </div>
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 700,
                fontSize: "22px",
                marginBottom: "8px",
              }}
            >
              Drop your résumé here
            </div>
            <p
              style={{
                fontSize: "15px",
                color: "#8A8175",
                margin: "0 0 24px",
              }}
            >
              PDF, DOCX, or TXT · up to 5MB
            </p>
            <span
              style={{
                display: "inline-block",
                background: "#C2683C",
                color: "#fff",
                fontSize: "15px",
                fontWeight: 600,
                padding: "13px 28px",
                borderRadius: "100px",
              }}
            >
              Choose file
            </span>
            <div
              style={{ marginTop: "18px", fontSize: "14px", color: "#8A8175" }}
            >
              or{" "}
              <span
                style={{
                  color: "#C2683C",
                  fontWeight: 600,
                  cursor: "pointer",
                }}
              >
                try with a sample résumé
              </span>
            </div>
          </div>

          {/* CHECKLIST */}
          <div
            style={{
              background: "#F3ECE2",
              border: "1px solid #EFE7DC",
              borderRadius: "20px",
              padding: "30px",
            }}
          >
            <div
              style={{
                fontFamily: "'Public Sans',sans-serif",
                fontWeight: 700,
                fontSize: "18px",
                marginBottom: "18px",
              }}
            >
              What we check
            </div>
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "15px",
              }}
            >
              {[
                "ATS-readable formatting",
                "Keyword & skill match",
                "Section structure",
                "Measurable impact",
                "Length & readability",
              ].map((item) => (
                <div
                  key={item}
                  style={{
                    display: "flex",
                    gap: "12px",
                    fontSize: "15px",
                    color: "#5B554C",
                  }}
                >
                  <span style={{ color: "#4F7C6A" }}>✓</span> {item}
                </div>
              ))}
            </div>
            <div
              style={{
                marginTop: "22px",
                paddingTop: "20px",
                borderTop: "1px solid #E2D9CC",
                fontSize: "13px",
                color: "#8A8175",
                lineHeight: 1.55,
              }}
            >
              Your file is analyzed privately and never shared with recruiters
              without your say-so.
            </div>
          </div>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
