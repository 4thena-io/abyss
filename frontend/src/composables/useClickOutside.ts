import { onMounted, onUnmounted, type Ref } from 'vue';

/**
 * Composable to detect clicks outside a target element
 * @param elementRef - Ref to the target element
 * @param callback - Function to call when a click outside is detected
 */
export function useClickOutside(
  elementRef: Ref<HTMLElement | null>,
  callback: () => void
) {
  const handleClick = (event: Event) => {
    if (elementRef.value && !elementRef.value.contains(event.target as Node)) {
      callback();
    }
  };

  onMounted(() => {
    document.addEventListener('click', handleClick);
  });

  onUnmounted(() => {
    document.removeEventListener('click', handleClick);
  });
}
