/**
 * Test setup file
 *
 * - Adds DOM matchers from @testing-library/jest-dom
 * - Provides minimal polyfills for browser APIs commonly used in
 *   components/tests (ResizeObserver, matchMedia) so tests run reliably
 *   in the JSDOM environment used by Vitest.
 */

import "@testing-library/jest-dom";

// Minimal ResizeObserver polyfill for the test environment.
// Tests in this repo don't need full ResizeObserver behavior, but some
// components guard for its existence and instantiate it — so providing
// a no-op implementation keeps tests stable.
if (typeof globalThis.ResizeObserver === "undefined") {
  class ResizeObserver {
    constructor(callback) {
      this.callback = callback;
      this.elements = new Set();
    }
    observe(element) {
      this.elements.add(element);
    }
    unobserve(element) {
      this.elements.delete(element);
    }
    disconnect() {
      this.elements.clear();
    }
  }
  globalThis.ResizeObserver = ResizeObserver;
}

// Minimal matchMedia polyfill. Some UI libraries or CSS-in-JS checks
// may call matchMedia; provide a no-op implementation.
if (typeof window.matchMedia === "undefined") {
  Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: (query) => ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {}, // deprecated API
      removeListener: () => {}, // deprecated API
      dispatchEvent: () => false,
    }),
  });
}
