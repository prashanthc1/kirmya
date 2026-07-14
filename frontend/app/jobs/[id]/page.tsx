"use client";

import React, { use } from "react";
import JobsPage from "../page";

interface PageProps {
  params: Promise<{ id: string }> | { id: string };
}

export default function JobDetailPage({ params }: PageProps) {
  // Next.js 15 uses promises for route parameters; safe lookup support for all versions:
  const resolvedParams = params && "then" in params ? use(params) : (params as { id: string });
  const jobId = resolvedParams?.id || null;

  return <JobsPage initialJobId={jobId} />;
}
