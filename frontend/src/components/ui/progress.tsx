import { cn } from "@/lib/utils";

export function Progress({ value, max = 100, className, color }: {
  value: number;
  max?: number;
  className?: string;
  color?: string;
}) {
  const percentage = Math.min((value / max) * 100, 100);

  return (
    <div className={cn("relative h-3 w-full overflow-hidden rounded-full bg-secondary", className)}>
      <div
        className={cn("h-full transition-all duration-500 ease-out rounded-full", color || "bg-primary")}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
}
