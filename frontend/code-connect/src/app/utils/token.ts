export const getToken = (): string | null => {
  try {
    if (!document.cookie) return null;

    const match = document.cookie.match(/sessionToken=([^;]+)/);
    return match ? match[1] : null;

  } catch (error) {
    console.error('Error getting session token:', error);
    return null;
  }
};