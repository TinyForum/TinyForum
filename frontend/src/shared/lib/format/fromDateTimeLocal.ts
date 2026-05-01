export const fromDateTimeLocal = (
  localValue: string | null | undefined,
): string | null => {
  if (!localValue) return null;
  try {
    const date = new Date(localValue);
    if (isNaN(date.getTime())) return null;
    return date.toISOString();
  } catch {
    return null;
  }
};
