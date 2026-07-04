import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function JobsPage() {
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
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Jobs" }]}
      />

      {/* HEADING */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "clamp(40px,5vw,64px) 40px clamp(24px,3vw,32px)",
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
          480 open roles
        </div>
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
          Roles that want your experience
        </h1>
        <p
          style={{
            fontSize: "clamp(16px,2vw,19px)",
            lineHeight: 1.6,
            color: "#5B554C",
            maxWidth: "560px",
            margin: 0,
          }}
        >
          Hand-vetted openings from recruiters who hire for proven track records.
        </p>
      </section>

      {/* SEARCH BAR */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "0 40px 28px",
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
              placeholder="Job title, skill, or company"
              style={{
                border: "none",
                outline: "none",
                background: "transparent",
                fontFamily: "'Public Sans',sans-serif",
                fontSize: "16px",
                color: "#2B2620",
                width: "100%",
                padding: "12px 0",
              }}
            />
          </div>
          <div
            style={{ width: "1px", height: "28px", background: "#EFE7DC" }}
          />
          <div
            style={{
              flex: "1 1 200px",
              display: "flex",
              alignItems: "center",
              gap: "10px",
              padding: "0 14px",
            }}
          >
            <span style={{ color: "#8A8175", fontSize: "16px" }}>⊙</span>
            <input
              placeholder="Location or remote"
              style={{
                border: "none",
                outline: "none",
                background: "transparent",
                fontFamily: "'Public Sans',sans-serif",
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
              fontFamily: "'Public Sans',sans-serif",
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
          <span
            style={{ fontSize: "13px", color: "#8A8175", alignSelf: "center" }}
          >
            Popular:
          </span>
          <button
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
            Remote
          </button>
          <button
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
            Operations
          </button>
          <button
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
            Engineering
          </button>
          <button
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
            Senior
          </button>
        </div>
      </section>

      {/* FILTER + RESULTS */}
      <section
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "0 40px clamp(56px,6vw,90px)",
          display: "grid",
          gridTemplateColumns: "260px 1fr",
          gap: "32px",
          alignItems: "start",
        }}
      >
        {/* FILTER SIDEBAR */}
        <aside
          style={{
            background: "#fff",
            border: "1px solid #EFE7DC",
            borderRadius: "18px",
            padding: "26px",
            position: "sticky",
            top: "120px",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              marginBottom: "20px",
            }}
          >
            <div
              style={{
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontWeight: 700,
                fontSize: "17px",
              }}
            >
              Filter
            </div>
            <button
              style={{
                border: "none",
                background: "none",
                cursor: "pointer",
                fontSize: "13px",
                color: "#C2683C",
                fontWeight: 600,
              }}
            >
              Clear
            </button>
          </div>
          <div style={{ marginBottom: "24px" }}>
            <div
              style={{
                fontSize: "13px",
                fontWeight: 700,
                letterSpacing: "0.06em",
                textTransform: "uppercase",
                color: "#8A8175",
                marginBottom: "12px",
              }}
            >
              Work type
            </div>
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "11px",
                fontSize: "15px",
                color: "#5B554C",
              }}
            >
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  cursor: "pointer",
                }}
              >
                <span
                  style={{
                    width: "18px",
                    height: "18px",
                    border: "1px solid #C9BEAD",
                    borderRadius: "5px",
                    display: "inline-block",
                  }}
                />
                Remote
              </label>
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  cursor: "pointer",
                }}
              >
                <span
                  style={{
                    width: "18px",
                    height: "18px",
                    border: "1px solid #C9BEAD",
                    borderRadius: "5px",
                    display: "inline-block",
                  }}
                />
                Hybrid
              </label>
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  cursor: "pointer",
                }}
              >
                <span
                  style={{
                    width: "18px",
                    height: "18px",
                    border: "1px solid #C9BEAD",
                    borderRadius: "5px",
                    display: "inline-block",
                  }}
                />
                On-site
              </label>
            </div>
          </div>
          <div>
            <div
              style={{
                fontSize: "13px",
                fontWeight: 700,
                letterSpacing: "0.06em",
                textTransform: "uppercase",
                color: "#8A8175",
                marginBottom: "12px",
              }}
            >
              Experience
            </div>
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "11px",
                fontSize: "15px",
                color: "#5B554C",
              }}
            >
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  cursor: "pointer",
                }}
              >
                <span
                  style={{
                    width: "18px",
                    height: "18px",
                    border: "1px solid #C9BEAD",
                    borderRadius: "5px",
                    display: "inline-block",
                  }}
                />
                Mid-level
              </label>
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  cursor: "pointer",
                }}
              >
                <span
                  style={{
                    width: "18px",
                    height: "18px",
                    border: "1px solid #C9BEAD",
                    borderRadius: "5px",
                    display: "inline-block",
                  }}
                />
                Senior
              </label>
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  cursor: "pointer",
                }}
              >
                <span
                  style={{
                    width: "18px",
                    height: "18px",
                    border: "1px solid #C9BEAD",
                    borderRadius: "5px",
                    display: "inline-block",
                  }}
                />
                Director+
              </label>
            </div>
          </div>
        </aside>

        {/* JOB CARDS */}
        <div>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              marginBottom: "18px",
              flexWrap: "wrap",
              gap: "12px",
            }}
          >
            <span style={{ fontSize: "15px", color: "#6B6357" }}>
              Showing <strong style={{ color: "#2B2620" }}>12</strong> of 480
              roles
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
              Sort: <strong style={{ color: "#2B2620" }}>Best match</strong> ▾
            </button>
          </div>
          <div
            style={{ display: "flex", flexDirection: "column", gap: "14px" }}
          >
            {/* JOB CARD 1 */}
            <div
              style={{
                position: "relative",
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "16px",
                padding: "24px",
                display: "flex",
                gap: "18px",
                alignItems: "flex-start",
              }}
            >
              <a
                href="/jobs/detail"
                style={{
                  flex: 1,
                  minWidth: 0,
                  display: "flex",
                  gap: "18px",
                  alignItems: "flex-start",
                }}
              >
                <span
                  style={{
                    flex: "none",
                    width: "52px",
                    height: "52px",
                    borderRadius: "13px",
                    background: "#F3ECE2",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "22px",
                    color: "#C2683C",
                  }}
                >
                  ◆
                </span>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "10px",
                      flexWrap: "wrap",
                      marginBottom: "6px",
                    }}
                  >
                    <h3
                      style={{
                        fontFamily: "'Bricolage Grotesque',sans-serif",
                        fontWeight: 700,
                        fontSize: "20px",
                        margin: 0,
                      }}
                    >
                      Director of Operations
                    </h3>
                    <span
                      style={{
                        fontSize: "12px",
                        color: "#4F7C6A",
                        background: "rgba(79,124,106,0.12)",
                        padding: "4px 10px",
                        borderRadius: "100px",
                        fontWeight: 600,
                      }}
                    >
                      Remote
                    </span>
                  </div>
                  <div
                    style={{
                      fontSize: "15px",
                      color: "#6B6357",
                      marginBottom: "14px",
                    }}
                  >
                    Northwind Logistics · Austin, TX (remote)
                  </div>
                  <div
                    style={{ display: "flex", gap: "8px", flexWrap: "wrap" }}
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
                      Team leadership
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
                      Process design
                    </span>
                  </div>
                </div>
              </a>
              <div
                style={{
                  flex: "none",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "flex-end",
                  gap: "10px",
                }}
              >
                <button
                  aria-label="Save job"
                  style={{
                    border: "1px solid #E2D9CC",
                    background: "#fff",
                    width: "38px",
                    height: "38px",
                    borderRadius: "10px",
                    cursor: "pointer",
                    color: "#8A8175",
                  }}
                >
                  ♡
                </button>
                <div style={{ textAlign: "right" }}>
                  <div
                    style={{
                      fontWeight: 700,
                      fontSize: "16px",
                      fontFamily: "'Bricolage Grotesque',sans-serif",
                    }}
                  >
                    $140k–$170k
                  </div>
                  <div
                    style={{
                      fontSize: "13px",
                      color: "#8A8175",
                      marginTop: "4px",
                    }}
                  >
                    2 days ago
                  </div>
                </div>
              </div>
            </div>

            {/* JOB CARD 2 */}
            <div
              style={{
                position: "relative",
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "16px",
                padding: "24px",
                display: "flex",
                gap: "18px",
                alignItems: "flex-start",
              }}
            >
              <a
                href="/jobs/detail"
                style={{
                  flex: 1,
                  minWidth: 0,
                  display: "flex",
                  gap: "18px",
                  alignItems: "flex-start",
                }}
              >
                <span
                  style={{
                    flex: "none",
                    width: "52px",
                    height: "52px",
                    borderRadius: "13px",
                    background: "#F3ECE2",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "22px",
                    color: "#4F7C6A",
                  }}
                >
                  ▣
                </span>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "10px",
                      flexWrap: "wrap",
                      marginBottom: "6px",
                    }}
                  >
                    <h3
                      style={{
                        fontFamily: "'Bricolage Grotesque',sans-serif",
                        fontWeight: 700,
                        fontSize: "20px",
                        margin: 0,
                      }}
                    >
                      Senior Product Manager
                    </h3>
                    <span
                      style={{
                        fontSize: "12px",
                        color: "#6A5FA0",
                        background: "rgba(106,95,160,0.12)",
                        padding: "4px 10px",
                        borderRadius: "100px",
                        fontWeight: 600,
                      }}
                    >
                      Hybrid
                    </span>
                  </div>
                  <div
                    style={{
                      fontSize: "15px",
                      color: "#6B6357",
                      marginBottom: "14px",
                    }}
                  >
                    Lumen Health · Boston, MA
                  </div>
                  <div
                    style={{ display: "flex", gap: "8px", flexWrap: "wrap" }}
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
                      Roadmapping
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
                      B2B SaaS
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
                      Analytics
                    </span>
                  </div>
                </div>
              </a>
              <div
                style={{
                  flex: "none",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "flex-end",
                  gap: "10px",
                }}
              >
                <button
                  aria-label="Save job"
                  style={{
                    border: "1px solid #E2D9CC",
                    background: "#fff",
                    width: "38px",
                    height: "38px",
                    borderRadius: "10px",
                    cursor: "pointer",
                    color: "#8A8175",
                  }}
                >
                  ♡
                </button>
                <div style={{ textAlign: "right" }}>
                  <div
                    style={{
                      fontWeight: 700,
                      fontSize: "16px",
                      fontFamily: "'Bricolage Grotesque',sans-serif",
                    }}
                  >
                    $150k–$185k
                  </div>
                  <div
                    style={{
                      fontSize: "13px",
                      color: "#8A8175",
                      marginTop: "4px",
                    }}
                  >
                    4 days ago
                  </div>
                </div>
              </div>
            </div>

            {/* JOB CARD 3 */}
            <div
              style={{
                position: "relative",
                background: "#fff",
                border: "1px solid #EFE7DC",
                borderRadius: "16px",
                padding: "24px",
                display: "flex",
                gap: "18px",
                alignItems: "flex-start",
              }}
            >
              <a
                href="/jobs/detail"
                style={{
                  flex: 1,
                  minWidth: 0,
                  display: "flex",
                  gap: "18px",
                  alignItems: "flex-start",
                }}
              >
                <span
                  style={{
                    flex: "none",
                    width: "52px",
                    height: "52px",
                    borderRadius: "13px",
                    background: "#F3ECE2",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "22px",
                    color: "#2B2620",
                  }}
                >
                  ●
                </span>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "10px",
                      flexWrap: "wrap",
                      marginBottom: "6px",
                    }}
                  >
                    <h3
                      style={{
                        fontFamily: "'Bricolage Grotesque',sans-serif",
                        fontWeight: 700,
                        fontSize: "20px",
                        margin: 0,
                      }}
                    >
                      People Operations Lead
                    </h3>
                    <span
                      style={{
                        fontSize: "12px",
                        color: "#4F7C6A",
                        background: "rgba(79,124,106,0.12)",
                        padding: "4px 10px",
                        borderRadius: "100px",
                        fontWeight: 600,
                      }}
                    >
                      Remote
                    </span>
                  </div>
                  <div
                    style={{
                      fontSize: "15px",
                      color: "#6B6357",
                      marginBottom: "14px",
                    }}
                  >
                    Atlas &amp; Co · Remote (US)
                  </div>
                  <div
                    style={{ display: "flex", gap: "8px", flexWrap: "wrap" }}
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
                      HR
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
                      Culture
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
                      Hiring
                    </span>
                  </div>
                </div>
              </a>
              <div
                style={{
                  flex: "none",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "flex-end",
                  gap: "10px",
                }}
              >
                <button
                  aria-label="Save job"
                  style={{
                    border: "1px solid #E2D9CC",
                    background: "#fff",
                    width: "38px",
                    height: "38px",
                    borderRadius: "10px",
                    cursor: "pointer",
                    color: "#8A8175",
                  }}
                >
                  ♡
                </button>
                <div style={{ textAlign: "right" }}>
                  <div
                    style={{
                      fontWeight: 700,
                      fontSize: "16px",
                      fontFamily: "'Bricolage Grotesque',sans-serif",
                    }}
                  >
                    $110k–$135k
                  </div>
                  <div
                    style={{
                      fontSize: "13px",
                      color: "#8A8175",
                      marginTop: "4px",
                    }}
                  >
                    1 week ago
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
