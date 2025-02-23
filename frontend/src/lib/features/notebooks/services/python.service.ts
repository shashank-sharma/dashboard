import { loadPyodide, type PyodideInterface } from "pyodide";

class PythonService {
    private pyodide: PyodideInterface | null = null;
    private namespace: any = {};

    async initialize() {
        if (!this.pyodide) {
            this.pyodide = await loadPyodide({
                indexURL: "https://cdn.jsdelivr.net/pyodide/v0.26.4/full/"
            });
            
            // Initialize the namespace with the capture function
            await this.pyodide.runPythonAsync(`
                import sys
                from io import StringIO

                def run_code(code):
                    old_stdout = sys.stdout
                    redirected_output = sys.stdout = StringIO()
                    try:
                        # Try to eval first (for expressions)
                        try:
                            result = eval(code, globals())
                            if result is not None:
                                print(repr(result))
                        except:
                            # If eval fails, try exec (for statements)
                            exec(code, globals())
                    finally:
                        sys.stdout = old_stdout
                        
                    return redirected_output.getvalue()
            `);
        }
        return this.pyodide;
    }

    async executeCode(code: string): Promise<string> {
        try {
            const pyodide = await this.initialize();
            const result = await pyodide.runPythonAsync(`run_code(${JSON.stringify(code)})`);
            return result;
        } catch (error) {
            throw error;
        }
    }
}

export const pythonService = new PythonService();