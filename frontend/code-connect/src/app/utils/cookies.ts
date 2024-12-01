export const getToken = (): string | null => {
  try {
    const cookies = document.cookie;
    if (!cookies) return null;

    const match = cookies.match(/session_token=([^;]+)/);
    return match ? match[1] : null;

  } catch (error) {
    console.error('Error getting session token:', error);
    return null;
  }
};

export const getUsername = (): string | null => {
  try {
    const cookies = document.cookie;
    if (!cookies) return null;

    const match = cookies.match(/username=([^;]+)/);
    return match ? match[1] : null;

  } catch (error) {
    console.error('Error getting username:', error);
    return null;
  }
}