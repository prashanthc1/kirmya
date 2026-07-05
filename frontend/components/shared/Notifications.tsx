"use client";

import React, { createContext, useContext, useState, useCallback } from "react";
import { Snackbar, Alert, AlertColor, Grow } from "@mui/material";

interface NotificationsContextType {
  showNotification: (message: string, severity?: AlertColor) => void;
}

const NotificationsContext = createContext<NotificationsContextType | undefined>(undefined);

export function useNotifications() {
  const context = useContext(NotificationsContext);
  if (!context) {
    throw new Error("useNotifications must be used within a NotificationsProvider");
  }
  return context;
}

interface NotificationsProviderProps {
  children: React.ReactNode;
}

export function NotificationsProvider({ children }: NotificationsProviderProps) {
  const [open, setOpen] = useState(false);
  const [message, setMessage] = useState("");
  const [severity, setSeverity] = useState<AlertColor>("success");

  const showNotification = useCallback((msg: string, sev: AlertColor = "success") => {
    setMessage(msg);
    setSeverity(sev);
    setOpen(true);
  }, []);

  const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }
    setOpen(false);
  };

  return (
    <NotificationsContext.Provider value={{ showNotification }}>
      {children}
      <Snackbar
        open={open}
        autoHideDuration={5000}
        onClose={handleClose}
        TransitionComponent={Grow}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
        sx={{
          zIndex: 2000,
        }}
      >
        <Alert
          onClose={handleClose}
          severity={severity}
          variant="filled"
          elevation={6}
          sx={{
            borderRadius: 3,
            fontWeight: 600,
            fontSize: "0.95rem",
            px: 2.5,
            py: 1,
            boxShadow: "0 10px 30px -5px rgba(43, 38, 32, 0.15)",
            "&.MuiAlert-filledSuccess": {
              backgroundColor: "secondary.main",
              color: "primary.contrastText",
            },
            "&.MuiAlert-filledError": {
              backgroundColor: "#ef4444",
            },
            "&.MuiAlert-filledWarning": {
              backgroundColor: "#f59e0b",
            },
            "&.MuiAlert-filledInfo": {
              backgroundColor: "#3b82f6",
            },
          }}
        >
          {message}
        </Alert>
      </Snackbar>
    </NotificationsContext.Provider>
  );
}
