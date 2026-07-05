/// Wails injects these at runtime in the desktop app.
declare global {
  interface Window {
    runtime: Record<string, (...args: unknown[]) => unknown>
    go: {
      main: {
        App: Record<string, (...args: unknown[]) => unknown>
      }
    }
  }
}

export {}
