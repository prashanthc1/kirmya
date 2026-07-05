import { createTheme } from "@mui/material/styles";

// Custom theme configuration for Kirmya.
// Aligns with the "Organic Premium Tech" styling language.
const theme = createTheme({
  palette: {
    primary: {
      main: "#D66838", // Sunset Coral (Comeback, energy)
      light: "#E38A63",
      dark: "#A64D24",
      contrastText: "#FFFFFF",
    },
    secondary: {
      main: "#37614D", // Eucalyptus Green (Growth, recovery)
      light: "#4E7B65",
      dark: "#234132",
      contrastText: "#FFFFFF",
    },
    background: {
      default: "#FCFAF7", // Alabaster Warm background
      paper: "#FFFFFF", // Pure Alabaster surfaces
    },
    text: {
      primary: "#2B2620", // Charcoal Bark
      secondary: "#6B6359", // Warm Slate
      disabled: "#A69E93",
    },
    divider: "rgba(43, 38, 32, 0.08)",
  },
  typography: {
    fontFamily: "var(--font-public-sans), sans-serif",
    h1: {
      fontFamily: "var(--font-bricolage), sans-serif",
      fontWeight: 800,
      letterSpacing: "-0.03em",
      lineHeight: 1.05,
    },
    h2: {
      fontFamily: "var(--font-bricolage), sans-serif",
      fontWeight: 800,
      letterSpacing: "-0.02em",
      lineHeight: 1.1,
    },
    h3: {
      fontFamily: "var(--font-bricolage), sans-serif",
      fontWeight: 800,
      letterSpacing: "-0.02em",
      lineHeight: 1.15,
    },
    h4: {
      fontFamily: "var(--font-bricolage), sans-serif",
      fontWeight: 700,
      letterSpacing: "-0.01em",
      lineHeight: 1.2,
    },
    h5: {
      fontFamily: "var(--font-bricolage), sans-serif",
      fontWeight: 700,
      letterSpacing: "-0.01em",
      lineHeight: 1.25,
    },
    h6: {
      fontFamily: "var(--font-bricolage), sans-serif",
      fontWeight: 700,
      lineHeight: 1.3,
    },
    body1: {
      fontSize: "1rem",
      lineHeight: 1.6,
      color: "#2B2620",
    },
    body2: {
      fontSize: "0.875rem",
      lineHeight: 1.6,
      color: "#6B6359",
    },
    button: {
      fontWeight: 600,
      textTransform: "none",
      letterSpacing: "0.01em",
    },
  },
  shape: {
    borderRadius: 12,
  },
  components: {
    MuiButton: {
      defaultProps: {
        disableElevation: true,
      },
      styleOverrides: {
        root: {
          borderRadius: 100, // Pill styling
          padding: "10px 24px",
          transition: "transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), background-color 0.2s ease, box-shadow 0.2s ease",
          "&:hover": {
            transform: "translateY(-1px)",
          },
          "&:active": {
            transform: "translateY(1px)",
          },
        },
        containedPrimary: {
          backgroundColor: "#D66838",
          color: "#FFFFFF",
          "&:hover": {
            backgroundColor: "#BE5729",
            boxShadow: "0 4px 12px rgba(214, 104, 56, 0.25)",
          },
        },
        outlinedPrimary: {
          borderColor: "rgba(214, 104, 56, 0.3)",
          color: "#D66838",
          "&:hover": {
            borderColor: "#D66838",
            backgroundColor: "rgba(214, 104, 56, 0.04)",
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 16,
          border: "1px solid rgba(43, 38, 32, 0.05)",
          boxShadow: "0 4px 20px -2px rgba(43, 38, 32, 0.02), 0 12px 40px -8px rgba(214, 104, 56, 0.03)",
          backgroundColor: "#FFFFFF",
          transition: "transform 0.4s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.4s ease, border-color 0.4s ease",
          overflow: "visible", // Allows internal glow filters to show
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          "& .MuiOutlinedInput-root": {
            borderRadius: 12,
            backgroundColor: "#FFFFFF",
            transition: "border-color 0.2s ease, box-shadow 0.2s ease",
            "& fieldset": {
              borderColor: "rgba(43, 38, 32, 0.12)",
            },
            "&:hover fieldset": {
              borderColor: "rgba(214, 104, 56, 0.3)",
            },
            "&.Mui-focused fieldset": {
              borderColor: "#D66838",
              borderWidth: "1.5px",
            },
          },
        },
      },
    },
  },
});

export default theme;
