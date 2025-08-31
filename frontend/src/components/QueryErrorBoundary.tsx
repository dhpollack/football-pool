import { useQueryErrorResetBoundary } from "@tanstack/react-query";
import ErrorBoundary from "./ErrorBoundary";
import type { ReactNode } from "react";

interface QueryErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
}

/**
 * Specialized error boundary for React Query errors.
 * Automatically resets React Query error state when retrying.
 */
export default function QueryErrorBoundary({
  children,
  fallback,
}: QueryErrorBoundaryProps) {
  const { reset } = useQueryErrorResetBoundary();

  return (
    <ErrorBoundary fallback={fallback} onReset={reset}>
      {children}
    </ErrorBoundary>
  );
}
