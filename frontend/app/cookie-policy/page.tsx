"use client";

import React from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

export default function CookiePolicyPage() {
  return (
    <div
      style={{
        background: "#FBF7F2",
        fontFamily: "'Public Sans', sans-serif",
        color: "#2B2620",
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Cookie Policy" }]} />

      <main style={{ flex: 1, maxWidth: "800px", margin: "0 auto", width: "100%", padding: "40px 24px 80px 24px" }}>
        <div style={{ fontSize: "13px", fontWeight: 700, letterSpacing: "0.12em", textTransform: "uppercase", color: "#C2683C", marginBottom: "10px" }}>
          Legal &amp; Compliance
        </div>
        <h1 style={{ fontWeight: 800, fontSize: "clamp(32px, 5vw, 48px)", lineHeight: 1.05, letterSpacing: "-0.03em", marginBottom: "28px" }}>
          Cookie Policy
        </h1>
        
        <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "20px", padding: "clamp(24px, 4vw, 36px)", display: "flex", flexDirection: "column", gap: "28px", lineHeight: "1.6" }}>
          
          <section>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>1. What Are Cookies</h2>
            <p style={{ margin: 0, color: "#5B554C" }}>
              Cookies are small text files stored on your computer or mobile device when you visit websites. They are widely used to make websites work, or work more efficiently, as well as to provide reporting information.
            </p>
          </section>

          <section>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>2. Why We Use Cookies</h2>
            <p style={{ margin: "0 0 12px 0", color: "#5B554C" }}>
              At Kirmya, we use cookies to provide a premium, seamless experience. Some cookies are required for technical reasons, while others enable personalization, analytics, and enhanced features.
            </p>
            <p style={{ margin: 0, color: "#5B554C" }}>
              We classify cookies into the following categories depending on their source and intent.
            </p>
          </section>

          <section>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>3. Cookie Categories &amp; Retention</h2>
            
            <div style={{ display: "flex", flexDirection: "column", gap: "16px", marginTop: "16px" }}>
              <div style={{ padding: "14px", border: "1px solid #EFE7DC", borderRadius: "12px", background: "#FCFAF7" }}>
                <strong style={{ display: "block", fontSize: "15px", color: "#2B2620" }}>Essential Cookies (Always Active)</strong>
                <span style={{ fontSize: "13px", color: "#5B554C" }}>
                  Used for sign-ins, session maintenance, CSRF double-submit protection, and cross-site scripting preventions. 
                  <br /><em>Retention:</em> Session duration or up to 30 days for &quot;Remember Me&quot;.
                </span>
              </div>

              <div style={{ padding: "14px", border: "1px solid #EFE7DC", borderRadius: "12px", background: "#FCFAF7" }}>
                <strong style={{ display: "block", fontSize: "15px", color: "#2B2620" }}>Functional Cookies</strong>
                <span style={{ fontSize: "13px", color: "#5B554C" }}>
                  Used to remember system language preferences, interface themes (light/dark mode), custom dashboard configurations, and workspace settings.
                  <br /><em>Retention:</em> 365 days.
                </span>
              </div>

              <div style={{ padding: "14px", border: "1px solid #EFE7DC", borderRadius: "12px", background: "#FCFAF7" }}>
                <strong style={{ display: "block", fontSize: "15px", color: "#2B2620" }}>Analytics &amp; Performance Cookies</strong>
                <span style={{ fontSize: "13px", color: "#5B554C" }}>
                  Logs anonymous visitor telemetry, response rates, API errors, and feature interactions to help our architects diagnose layout lags.
                  <br /><em>Retention:</em> 180 days.
                </span>
              </div>

              <div style={{ padding: "14px", border: "1px solid #EFE7DC", borderRadius: "12px", background: "#FCFAF7" }}>
                <strong style={{ display: "block", fontSize: "15px", color: "#2B2620" }}>Marketing &amp; Sponsorship Cookies</strong>
                <span style={{ fontSize: "13px", color: "#5B554C" }}>
                  Kept for future advertising integrations. Disabled by default currently since the platform is 100% free with no paid listings.
                  <br /><em>Retention:</em> 365 days.
                </span>
              </div>
            </div>
          </section>

          <section>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>4. Third-Party Cookies</h2>
            <p style={{ margin: 0, color: "#5B554C" }}>
              In addition to our first-party cookies, we may integrate third-party tools (such as anonymous analytics trackers or performance loggers). These third parties can place cookies on your device, but they are blocked from running unless you explicitly enable functional/analytics cookie categories.
            </p>
          </section>

          <section>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>5. Browser Management &amp; Deletion</h2>
            <p style={{ margin: "0 0 12px 0", color: "#5B554C" }}>
              You can instruct your web browser to refuse or delete cookies. If you choose to refuse essential cookies, please note that you will not be able to log in or maintain an active career account workspace.
            </p>
            <div style={{ background: "#FCFAF7", border: "1px dashed #E2D9CC", borderRadius: "12px", padding: "14px", fontSize: "13px", color: "#8A8175" }}>
              <strong>To clear cookies in Chrome/Firefox/Safari:</strong> Go to Settings → Privacy &amp; Security → Cookies and Site Data → Manage or Delete Cookies.
            </div>
          </section>

          <section>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>6. Policy Updates</h2>
            <p style={{ margin: 0, color: "#5B554C" }}>
              We may update this Cookie Policy from time to time. When we make substantial changes, we will update the consent version, which will automatically prompt you to review and confirm your preferences on your next visit.
            </p>
          </section>

          <section style={{ borderTop: "1px solid #EFE7DC", paddingTop: "20px" }}>
            <h2 style={{ fontSize: "20px", fontWeight: 700, margin: "0 0 12px 0" }}>7. Contact Information</h2>
            <p style={{ margin: 0, color: "#5B554C" }}>
              If you have any questions about our cookie usage, please reach out to our privacy engineers at{" "}
              <a href="mailto:privacy@kirmya.com" style={{ color: "#C2683C", fontWeight: 600, textDecoration: "none" }}>
                privacy@kirmya.com
              </a>.
            </p>
          </section>

        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
