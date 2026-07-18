"use client";

import React from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Typography,
  Slide,
} from "@mui/material";
import { TransitionProps } from "@mui/material/transitions";
import CloseIcon from "@mui/icons-material/Close";

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

interface MuiModalProps {
  open: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  actions?: React.ReactNode;
  maxWidth?: "xs" | "sm" | "md" | "lg" | "xl";
  fullWidth?: boolean;
}

export default function MuiModal({
  open,
  onClose,
  title,
  children,
  actions,
  maxWidth = "sm",
  fullWidth = true,
}: MuiModalProps) {
  return (
    <Dialog
      open={open}
      TransitionComponent={Transition}
      keepMounted
      onClose={onClose}
      maxWidth={maxWidth}
      fullWidth={fullWidth}
      aria-labelledby="mui-modal-title"
      aria-describedby="mui-modal-description"
      PaperProps={{
        elevation: 0,
        sx: {
          borderRadius: 4,
          border: "1px solid var(--border)",
          boxShadow: "0 24px 64px -12px rgba(0, 0, 0, 0.16)",
          padding: 1.5,
          background: "var(--card)",
        },
      }}
    >
      {/* Modal Title with Close Button */}
      <DialogTitle
        id="mui-modal-title"
        sx={{
          m: 0,
          p: 2,
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <Typography
          variant="h5"
          component="span"
          sx={{
            fontFamily: "var(--font-public-sans), sans-serif",
            fontWeight: 800,
            color: "text.primary",
            letterSpacing: "-0.01em",
          }}
        >
          {title}
        </Typography>
        <IconButton
          aria-label="close modal"
          onClick={onClose}
          sx={{
            color: "text.secondary",
            "&:hover": {
              backgroundColor: "rgba(43, 38, 32, 0.05)",
            },
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      {/* Modal Content */}
      <DialogContent
        id="mui-modal-description"
        sx={{
          p: 2,
          color: "text.primary",
        }}
      >
        {children}
      </DialogContent>

      {/* Modal Actions */}
      {actions && (
        <DialogActions
          sx={{
            p: 2,
            gap: 1.5,
            justifyContent: "flex-end",
          }}
        >
          {actions}
        </DialogActions>
      )}
    </Dialog>
  );
}
