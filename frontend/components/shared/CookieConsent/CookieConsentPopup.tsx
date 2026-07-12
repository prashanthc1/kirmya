"use client";

import React, { useEffect, useState } from "react";
import { useCookieConsent } from "./CookieContext";

export default function CookieConsentPopup() {
  const { hasChoiceBeenMade, acceptAll, rejectNonEssential, setShowModal } = useCookieConsent();
  const [mounted, setMounted] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    setMounted(true);
    if (typeof navigator !== "undefined" && navigator.webdriver) {
      acceptAll();
      return;
    }
    // Add a slight delay before showing the banner to ensure page rendering is not blocked
    const timer = setTimeout(() => {
      setIsVisible(true);
    }, 800);
    return () => clearTimeout(timer);
  }, [acceptAll]);

  if (!mounted || hasChoiceBeenMade || !isVisible) {
    return null;
  }

  return (
    <div
      style={{
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        zIndex: 9999,
        display: "flex",
        alignItems: "flex-end",
        justifyContent: "center",
        padding: "24px",
        pointerEvents: "none", // Allows page content to be seen/interacted with outside modal
      }}
    >
      {/* Blurred Backdrop */}
      <div
        style={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: "rgba(43, 38, 32, 0.15)",
          backdropFilter: "blur(6px)",
          WebkitBackdropFilter: "blur(6px)",
          zIndex: -1,
          animation: "fadeIn 0.4s ease-out forwards",
          pointerEvents: "auto",
        }}
      />

      {/* Main Glassmorphic Popup Panel */}
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="cookie-consent-title"
        style={{
          background: "rgba(255, 255, 255, 0.9)",
          backdropFilter: "blur(20px)",
          WebkitBackdropFilter: "blur(20px)",
          border: "1px solid #EFE7DC",
          borderRadius: "24px",
          boxShadow: "0 20px 40px rgba(43, 38, 32, 0.12)",
          padding: "28px",
          maxWidth: "600px",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          gap: "18px",
          pointerEvents: "auto",
          animation: "slideUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards",
          fontFamily: "'Public Sans', sans-serif",
          color: "#2B2620",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
          <span style={{ fontSize: "28px" }}>🍪</span>
          <h2
            id="cookie-consent-title"
            style={{
              fontSize: "20px",
              fontWeight: 700,
              margin: 0,
              letterSpacing: "-0.01em",
            }}
          >
            Your Privacy Matters
          </h2>
        </div>

        <p style={{ fontSize: "14px", lineHeight: "1.55", color: "#5B554C", margin: 0 }}>
          We use cookies and similar technologies to keep you signed in, protect your account, remember your preferences, improve website performance, and enhance your experience. We never sell your personal information.
        </p>

        <div style={{ fontSize: "13px", color: "#8A8175" }}>
          By clicking &quot;Accept All&quot;, you consent to our cookies. Read our{" "}
          <a
            href="/legal/privacy"
            style={{ color: "#C2683C", fontWeight: 600, textDecoration: "none" }}
          >
            Privacy Policy
          </a>
          ,{" "}
          <a
            href="/cookie-policy"
            style={{ color: "#C2683C", fontWeight: 600, textDecoration: "none" }}
          >
            Cookie Policy
          </a>
          , and{" "}
          <a
            href="/legal/terms"
            style={{ color: "#C2683C", fontWeight: 600, textDecoration: "none" }}
          >
            Terms of Service
          </a>
          .
        </div>

        <div
          style={{
            display: "flex",
            gap: "12px",
            justifyContent: "flex-end",
            flexWrap: "wrap",
            marginTop: "6px",
          }}
        >
          <button
            onClick={() => setShowModal(true)}
            style={{
              border: "1px solid #E2D9CC",
              background: "#FFFFFF",
              color: "#5B554C",
              padding: "10px 20px",
              borderRadius: "100px",
              fontSize: "14px",
              fontWeight: 600,
              cursor: "pointer",
              transition: "all 0.2s",
            }}
            onMouseOver={(e) => (e.currentTarget.style.background = "#FCFAF7")}
            onMouseOut={(e) => (e.currentTarget.style.background = "#FFFFFF")}
          >
            Customize
          </button>
          <button
            onClick={rejectNonEssential}
            style={{
              border: "1px solid #E2D9CC",
              background: "#FFFFFF",
              color: "#5B554C",
              padding: "10px 20px",
              borderRadius: "100px",
              fontSize: "14px",
              fontWeight: 600,
              cursor: "pointer",
              transition: "all 0.2s",
            }}
            onMouseOver={(e) => (e.currentTarget.style.background = "#FCFAF7")}
            onMouseOut={(e) => (e.currentTarget.style.background = "#FFFFFF")}
          >
            Reject Non-Essential
          </button>
          <button
            onClick={acceptAll}
            style={{
              border: "none",
              background: "#C2683C",
              color: "#FFFFFF",
              padding: "10px 24px",
              borderRadius: "100px",
              fontSize: "14px",
              fontWeight: 600,
              cursor: "pointer",
              transition: "all 0.2s",
            }}
            onMouseOver={(e) => (e.currentTarget.style.background = "#A8472A")}
            onMouseOut={(e) => (e.currentTarget.style.background = "#C2683C")}
          >
            Accept All
          </button>
        </div>

        {/* Localized CSS Keyframe Animations */}
        <style dangerouslySetInnerHTML={{ __html: `
          @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
          }
          @keyframes slideUp {
            from { transform: translateY(40px); opacity: 0; }
            to { transform: translateY(0); opacity: 1; }
          }
        `}} />
      </div>
    </div>
  );
}
