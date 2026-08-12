/**
 * Composable for common formatting utilities
 */
export function useFormatters() {
  /**
   * Format duration in seconds to human-readable string
   * @param seconds - Duration in seconds
   * @returns Formatted string like "2m 30s" or "-" if falsy
   */
  const formatDuration = (seconds: number): string => {
    if (!seconds) return '-';
    if (seconds < 60) return `${seconds}s`;
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}m ${secs}s`;
  };

  /**
   * Format a date string to relative time (e.g., "5m ago", "2h ago")
   * @param dateString - ISO date string
   * @returns Relative time string or "-" if falsy
   */
  const formatRelativeTime = (dateString: string): string => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const mins = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (mins < 1) return 'just now';
    if (mins < 60) return `${mins}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;
    return date.toLocaleDateString();
  };

  /**
   * Format a URL to ensure it has a protocol
   * @param url - URL string that may or may not have a protocol
   * @returns URL with https:// prefix if no protocol present
   */
  const formatUrl = (url: string): string => {
    if (!url) return '#';
    if (url.startsWith('http://') || url.startsWith('https://')) return url;
    return `https://${url}`;
  };

  return {
    formatDuration,
    formatRelativeTime,
    formatUrl,
  };
}
