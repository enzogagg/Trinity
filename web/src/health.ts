export type Health = { status: string };

export const isHealthy = (h: Health): boolean => h.status === "ok";
