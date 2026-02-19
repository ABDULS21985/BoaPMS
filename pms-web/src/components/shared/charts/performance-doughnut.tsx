"use client";

import { PieChart, Pie, Cell, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CHART_COLORS } from "@/lib/scorecard-helpers";

interface PerformanceDoughnutProps {
  title: string;
  percentage: number;
  earnedColor?: string;
  remainingColor?: string;
  className?: string;
}

export function PerformanceDoughnut({
  title,
  percentage,
  earnedColor = CHART_COLORS.earned,
  remainingColor = CHART_COLORS.remaining,
  className,
}: PerformanceDoughnutProps) {
  const data = [
    { name: "Earned", value: percentage },
    { name: "Remaining", value: Math.max(0, 100 - percentage) },
  ];

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="relative flex items-center justify-center" style={{ height: 250 }}>
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie data={data} cx="50%" cy="50%" innerRadius={60} outerRadius={90} dataKey="value" startAngle={90} endAngle={-270}>
                <Cell fill={earnedColor} />
                <Cell fill={remainingColor} />
              </Pie>
            </PieChart>
          </ResponsiveContainer>
          <div className="absolute text-3xl font-bold">{percentage.toFixed(1)}%</div>
        </div>
      </CardContent>
    </Card>
  );
}
