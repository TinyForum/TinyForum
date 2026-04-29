interface StatCardProps {
  title: string;
  value: number;
  icon: string;
  color: string;
}

export function StatCard({ title, value, icon, color }: StatCardProps) {
  return (
    <div className="card bg-base-100 border border-base-300 hover:shadow-lg transition-shadow">
      <div className="card-body">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-base-content/60 text-sm">{title}</p>
            <p className="text-3xl font-bold mt-1">{value}</p>
          </div>
          <div className={`text-3xl ${color}`}>{icon}</div>
        </div>
      </div>
    </div>
  );
}
