import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { useFormatters } from './useFormatters';

describe('useFormatters', () => {
  const { formatDuration, formatRelativeTime, formatUrl } = useFormatters();

  describe('formatDuration', () => {
    it('returns "-" for falsy input', () => {
      expect(formatDuration(0)).toBe('-');
    });

    it('formats sub-minute durations as seconds', () => {
      expect(formatDuration(45)).toBe('45s');
    });

    it('formats durations over a minute as minutes and seconds', () => {
      expect(formatDuration(150)).toBe('2m 30s');
    });
  });

  describe('formatRelativeTime', () => {
    beforeEach(() => {
      vi.useFakeTimers();
      vi.setSystemTime(new Date('2026-07-18T12:00:00Z'));
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it('returns "-" for a falsy date string', () => {
      expect(formatRelativeTime('')).toBe('-');
    });

    it('returns "just now" for the current moment', () => {
      expect(formatRelativeTime('2026-07-18T11:59:50Z')).toBe('just now');
    });

    it('formats minutes ago', () => {
      expect(formatRelativeTime('2026-07-18T11:45:00Z')).toBe('15m ago');
    });

    it('formats hours ago', () => {
      expect(formatRelativeTime('2026-07-18T09:00:00Z')).toBe('3h ago');
    });

    it('formats days ago', () => {
      expect(formatRelativeTime('2026-07-16T12:00:00Z')).toBe('2d ago');
    });

    it('falls back to a locale date string beyond a week', () => {
      const result = formatRelativeTime('2026-07-01T12:00:00Z');
      expect(result).not.toMatch(/ago$/);
    });
  });

  describe('formatUrl', () => {
    it('returns "#" for a falsy url', () => {
      expect(formatUrl('')).toBe('#');
    });

    it('leaves URLs with an existing protocol untouched', () => {
      expect(formatUrl('https://example.com')).toBe('https://example.com');
      expect(formatUrl('http://example.com')).toBe('http://example.com');
    });

    it('prefixes bare hosts with https://', () => {
      expect(formatUrl('example.com')).toBe('https://example.com');
    });
  });
});
