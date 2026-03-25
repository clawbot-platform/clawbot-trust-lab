type StatusTone = "good" | "warn" | "bad" | "neutral";

const toneMap: Record<string, StatusTone> = {
  improved: "good",
  clean: "good",
  accepted: "good",
  mixed: "neutral",
  suspicious: "warn",
  step_up_required: "warn",
  needs_follow_up: "warn",
  regressed: "bad",
  blocked: "bad",
  false_signal: "bad",
  new_blind_spot_discovered: "bad",
  duplicate: "neutral"
};

export function StatusPill({ value }: { value: string }) {
  const tone = toneMap[value] ?? "neutral";
  return <span className={`pill pill-${tone}`}>{value.replace(/_/g, " ")}</span>;
}
