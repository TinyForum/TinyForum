// ===================== 自定义 Hook =====================
export interface PasswordStrength {
  score: number;
  level: string;
  color: string;
  message: string;
}

export function usePasswordStrength(password: string): PasswordStrength {
  if (!password) {
    return { score: 0, level: "无", color: "bg-gray-200", message: "" };
  }

  let score = 0;
  if (password.length >= 6) score++;
  if (password.length >= 10) score++;
  if (/[a-z]/.test(password)) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^A-Za-z0-9]/.test(password)) score++;

  if (score <= 2) {
    return {
      score,
      level: "弱",
      color: "bg-red-500",
      message: "密码强度较弱，建议增加复杂度",
    };
  }
  if (score <= 4) {
    return {
      score,
      level: "中",
      color: "bg-yellow-500",
      message: "密码强度中等，可以更强一些",
    };
  }
  return {
    score,
    level: "强",
    color: "bg-green-500",
    message: "密码强度很好！",
  };
}
