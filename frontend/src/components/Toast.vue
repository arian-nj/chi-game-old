Heavilly AI Generated

<script lang="ts">
import { reactive, readonly, computed } from 'vue';

// --- TYPE DEFINITIONS ---
type ToastType = 'success' | 'error' | 'info' | 'warning';
type ToastPosition = 'top-center' | 'top-right' | 'top-left' | 'bottom-center' | 'bottom-right' | 'bottom-left';

interface Toast {
  id: number;
  message: string;
  type: ToastType;
  duration: number;
  // Internal properties for pause/resume logic
  _timeoutId?: ReturnType<typeof setTimeout>;
  _startTime?: number;
  _remaining?: number;
}

interface ToastOptions {
  type?: ToastType;
  duration?: number;
}

interface ToastConfig {
  duration: number;
  position: ToastPosition;
  pauseOnHover: boolean;
  limit: number;
}

// A function that can be called to show a toast, and also has methods for specific toast types.
type ToastApi = {
  (message: string, options?: Omit<ToastOptions, 'type'>): void;
  success: (message: string, options?: Omit<ToastOptions, 'type'>) => void;
  error: (message: string, options?: Omit<ToastOptions, 'type'>) => void;
  info: (message: string, options?: Omit<ToastOptions, 'type'>) => void;
  warning: (message: string, options?: Omit<ToastOptions, 'type'>) => void;
}

interface UseToastReturn {
  toast: ToastApi;
  remove: (id: number) => void;
  toasts: Readonly<Toast[]>;
  config: Readonly<ToastConfig>;
  // Internal functions for the component to use
  _pause: (toast: Toast) => void;
  _resume: (toast: Toast) => void;
}

// --- STATE MANAGEMENT (SINGLETON) ---
const toasts = reactive<Toast[]>([]);

// Global configuration with defaults
const toastConfig = reactive<ToastConfig>({
  duration: 4000,
  position: 'top-center',
  pauseOnHover: true,
  limit: 5, // New setting for the maximum number of toasts
});

/**
 * Configure the global toast settings.
 * @param options - The configuration options to set.
 */
export function configureToast(options: Partial<ToastConfig>) {
  Object.assign(toastConfig, options);
}

// --- COMPOSABLE LOGIC ---
/**
 * A composable to manage and display toast notifications.
 */
export function useToast(): UseToastReturn {
  const remove = (id: number) => {
    const index = toasts.findIndex((t) => t.id === id);
    if (index !== -1) {
      const toast = toasts[index];
      if (toast._timeoutId) clearTimeout(toast._timeoutId);
      toasts.splice(index, 1);
    }
  };

  const pause = (toast: Toast) => {
    if (toast._startTime === undefined || toast._remaining === undefined) return;
    clearTimeout(toast._timeoutId);
    toast._remaining -= (Date.now() - toast._startTime);
  };

  const resume = (toast: Toast) => {
    if (toast._remaining === undefined) return;
    toast._startTime = Date.now();
    toast._timeoutId = setTimeout(() => {
      remove(toast.id);
    }, toast._remaining);
  };

  const show = (message: string, options: ToastOptions = {}) => {
    const id = Date.now() + Math.random();
    const duration = options.duration || toastConfig.duration;

    const toast: Toast = reactive({
      id,
      message,
      type: options.type || 'info',
      duration: duration,
      _startTime: Date.now(),
      _remaining: duration,
      _timeoutId: setTimeout(() => remove(id), duration),
    });

    toasts.unshift(toast);

    // If the number of toasts now exceeds the limit, remove the oldest one.
    if (toasts.length > toastConfig.limit) {
      const oldestToast = toasts.pop(); // .pop() removes the last (oldest) item
      if (oldestToast?._timeoutId) {
        clearTimeout(oldestToast._timeoutId);
      }
    }
  };

  const toast: ToastApi = (message: string, options: Omit<ToastOptions, 'type'> = {}) => {
    show(message, { ...options, type: 'info' });
  };

  toast.success = (message, options) => show(message, { ...options, type: 'success' });
  toast.error = (message, options) => show(message, { ...options, type: 'error' });
  toast.info = (message, options) => show(message, { ...options, type: 'info' });
  toast.warning = (message, options) => show(message, { ...options, type: 'warning' });

  return {
    toast,
    remove,
    toasts: readonly(toasts),
    config: readonly(toastConfig),
    _pause: pause,
    _resume: resume,
  };
}
</script>

<script setup lang="ts">
import gsap from 'gsap';

// Use the composable within the component itself to get access to the state.
const { toasts, remove, config, _pause, _resume } = useToast();

const handleMouseEnter = (toast: Toast) => {
  if (config.pauseOnHover) {
    _pause(toast);
  }
};

const handleMouseLeave = (toast: Toast) => {
  if (config.pauseOnHover) {
    _resume(toast);
  }
};


// --- ANIMATION HOOKS ---
const onEnter = (el: Element, done: () => void) => {
  // Find the icon element within the toast
  const icon = el.querySelector('[data-toast-icon]');

  // Animate the main toast container
  gsap.fromTo(
    el,
    { y: -60, opacity: 0, scale: 0.9 },
    { y: 0, opacity: 1, scale: 1, duration: 0.4, ease: 'power3.out', onComplete: done }
  );

  // Animate the icon with a pop effect after the container starts appearing
  if (icon) {
    gsap.fromTo(
      icon,
      { scale: 0, opacity: 0, rotation: -30 },
      {
        scale: 1,
        opacity: 1,
        rotation: 0,
        duration: 0.4,
        ease: 'back.out(1.7)',
        delay: 0.15 // Start a little after the main animation
      }
    );
  }
};

const onLeave = (el: Element, done: () => void) => {
  gsap.to(el, {
    opacity: 0,
    y: -20,
    scale: 0.95,
    duration: 0.3,
    ease: 'power2.in',
    onComplete: done,
  });
};

// --- DYNAMIC STYLING ---
const containerClasses = computed(() => {
  const baseClasses = 'fixed z-[100] w-full max-w-sm sm:max-w-md px-4';
  const positionClasses: { [key in ToastPosition]: string } = {
    'top-center': 'top-5 left-1/2 -translate-x-1/2',
    'top-right': 'top-5 right-5',
    'top-left': 'top-5 left-5',
    'bottom-center': 'bottom-5 left-1/2 -translate-x-1/2',
    'bottom-right': 'bottom-5 right-5',
    'bottom-left': 'bottom-5 left-5',
  };
  return [baseClasses, positionClasses[config.position]];
});

const toastClasses = computed(() => (type: ToastType) => {
  const baseClasses = 'flex items-center w-full p-4 rounded-lg shadow-xl space-x-4 bg-gray-800 text-gray-100 border-l-4';
  const typeClasses = {
    success: 'border-green-500',
    error: 'border-red-500',
    warning: 'border-yellow-500',
    info: 'border-blue-500',
  };
  return [baseClasses, typeClasses[type]];
});

const iconClasses = computed(() => (type: ToastType) => {
  const baseClasses = 'inline-flex items-center justify-center flex-shrink-0 w-8 h-8';
  const typeClasses = {
    success: 'text-green-400',
    error: 'text-red-400',
    warning: 'text-yellow-400',
    info: 'text-blue-400',
  };
  return [baseClasses, typeClasses[type]];
});

</script>

<template>
  <div :class="containerClasses" aria-live="assertive">
    <TransitionGroup tag="div" class="flex flex-col items-center justify-center space-y-4" @enter="onEnter"
      @leave="onLeave">
      <div v-for="toast in toasts" :key="toast.id" :class="toastClasses(toast.type)" role="alert"
        @mouseenter="handleMouseEnter(toast)" @mouseleave="handleMouseLeave(toast)">
        <!-- Icon based on type -->
        <div :class="iconClasses(toast.type)" data-toast-icon>
          <!-- Success Icon -->
          <svg v-if="toast.type === 'success'" class="w-7 h-7" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <!-- Error Icon -->
          <svg v-if="toast.type === 'error'" class="w-7 h-7" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
              clip-rule="evenodd"></path>
          </svg>
          <!-- Warning Icon -->
          <svg v-if="toast.type === 'warning'" class="w-7 h-7" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.21 3.03-1.742 3.03H4.42c-1.532 0-2.492-1.696-1.742-3.03l5.58-9.92zM10 13a1 1 0 110-2 1 1 0 010 2zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
              clip-rule="evenodd"></path>
          </svg>
          <!-- Info Icon -->
          <svg v-if="toast.type === 'info'" class="w-7 h-7" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
              clip-rule="evenodd"></path>
          </svg>
        </div>

        <div class="text-base font-medium flex-grow mr-2">{{ toast.message }}</div>

        <!-- Close Button -->
        <button type="button"
          class="ml-auto flex-shrink-0 p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white inline-flex items-center justify-center h-9 w-9 text-gray-400 hover:text-white hover:bg-gray-700"
          :aria-label="`Close toast for notification: ${toast.message}`" @click="remove(toast.id)">
          <span class="sr-only">Close</span>
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
              clip-rule="evenodd"></path>
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
