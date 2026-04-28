export const StatItem = ({
  label,
  value,
}: {
  label: string;
  value: number;
}) => (
  <div className="flex justify-between">
    <span className="text-muted-foreground">{label}</span>
    <span className="font-medium">{value.toLocaleString()}</span>
  </div>
);
