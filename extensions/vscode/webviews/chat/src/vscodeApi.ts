// This file provides a simple, type-safe wrapper for
// communicating with the VS Code extension.

/**
 * A type-safe wrapper for the VS Code API.
 */
class VSCodeAPI {
  // We specify 'any' here because the VS Code API is not
  // available in a standard browser environment.
  // We'll provide a mock for development if needed.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private vsCodeApi: any;

  constructor() {
    // Check if the 'acquireVsCodeApi' function exists
    if (typeof acquireVsCodeApi === 'function') {
      this.vsCodeApi = acquireVsCodeApi();
    } else {
      // Provide a mock API for local development
      // in a browser
      console.warn('VS Code API not found. Using mock implementation.');
      this.vsCodeApi = {
        postMessage: (message: unknown) => {
          console.log('Mock postMessage:', message);
        },
        // You can add more mock methods as needed
      };
    }
  }

  public api() {
    return this.vsCodeApi
  }

  public getState() {
    return this.vsCodeApi.getState()
  }

  public setState(val: any) {
    this.vsCodeApi.setState(val)
  }


  /**
   * Posts a message to the VS Code extension.
   * @param message The message to send.
   */
  public postMessage(message: unknown) {
    this.vsCodeApi.postMessage(message);
  }

  /**
   * Adds a listener for messages from the VS Code extension.
   * @param callback The callback to run when a message is received.
   * @returns A function to remove the event listener.
   */
  public onMessage(callback: (event: MessageEvent) => void): () => void {
    const handler = (event: MessageEvent) => {
      callback(event);
    };
    window.addEventListener('message', handler);
    return () => window.removeEventListener('message', handler);
  }
}

// Create and export a single instance of the API
export const vscodeApi = new VSCodeAPI();