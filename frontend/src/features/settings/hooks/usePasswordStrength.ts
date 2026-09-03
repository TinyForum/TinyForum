// ===================== 自定义 Hook =====================
// 返回翻译 key，消费者需用 t() 转换
export interface PasswordStrength {
  score: number;
  level: string;
  color: string;
  message: string;
}

export function usePasswordStrength(password: string): PasswordStrength {
  if (!password) {
    return { score: 0, level: "password_strength_none", color: "bg-gray-200", message: "" };
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
      level: "password_strength_weak",
      color: "bg-red-500",
      message: "password_strength_weak_message",
    };
  }
  if (score <= 4) {
    return {
      score,
      level: "password_strength_medium",
      color: "bg-yellow-500",
      message: "password_strength_medium_message",
    };
  }
  return {
    score,
    level: "password_strength_strong",
    color: "bg-green-500",
    message: "password_strength_strong_message",
  };
}
