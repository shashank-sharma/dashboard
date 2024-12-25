declare module 'pyodide' {
    export interface PyodideInterface {
      loadPackage(packages: string | string[]): Promise<void>;
      runPythonAsync(code: string, options?: { globals?: any }): Promise<any>;
      setStderr(options: { write: (text: string) => void; flush: () => void }): void;
      setStdout(options: { write: (text: string) => void; flush: () => void }): void;
      globals: {
        get(name: string): any;
      };
    }
  
    export function loadPyodide(options?: {
      indexURL?: string;
    }): Promise<PyodideInterface>;
  }