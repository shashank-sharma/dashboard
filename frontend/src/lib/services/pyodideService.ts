import type { PyodideInterface } from 'pyodide';

let pyodide: PyodideInterface | null = null;
let isLoading = false;
let loadPromise: Promise<void> | null = null;

export const pyodideService = {
    async initialize() {
        if (pyodide) return pyodide;
        if (isLoading) return loadPromise;

        isLoading = true;
        loadPromise = new Promise(async (resolve, reject) => {
            try {
                // Load pyodide.js script
                const script = document.createElement('script');
                script.src = 'https://cdn.jsdelivr.net/pyodide/v0.24.1/full/pyodide.js';
                document.head.appendChild(script);

                await new Promise((resolve) => script.onload = resolve);

                // Initialize Pyodide
                // @ts-ignore - Pyodide will be available globally
                pyodide = await loadPyodide({
                    indexURL: 'https://cdn.jsdelivr.net/pyodide/v0.24.1/full/',
                    stdout: (text: string) => {
                        console.log("Python stdout:", text);
                    },
                    stderr: (text: string) => {
                        console.error("Python stderr:", text);
                    }
                });

                // Install packages
                await pyodide.loadPackage(['numpy', 'pandas']);
                
                // Set up Python environment with proper I/O handling
                await pyodide.runPythonAsync(`
                    import sys
                    import io

                    class StringIO(io.StringIO):
                        def __init__(self):
                            super().__init__()
                            self.output = []

                        def write(self, text):
                            self.output.append(text)
                            return len(text)

                        def getvalue(self):
                            return ''.join(self.output)

                    sys.stdout = StringIO()
                    sys.stderr = StringIO()
                `);

                isLoading = false;
                resolve();
            } catch (error) {
                isLoading = false;
                console.error("Error initializing Pyodide:", error);
                reject(error);
            }
        });

        return loadPromise;
    },

    async executeCode(code: string): Promise<{ output: string; error: string | null }> {
        if (!pyodide) {
            await this.initialize();
        }

        try {
            // Reset stdout and stderr before each execution
            await pyodide.runPythonAsync(`
                sys.stdout = StringIO()
                sys.stderr = StringIO()
            `);

            // Execute the code
            const result = await pyodide.runPythonAsync(`
                try:
                    __result = eval(${JSON.stringify(code)})
                    if __result is not None:
                        print(repr(__result))
                except SyntaxError:
                    exec(${JSON.stringify(code)})
                
                output = sys.stdout.getvalue().strip()
                error = sys.stderr.getvalue().strip()
                (output, error)
            `);

            const [output, error] = result;
            
            return {
                output: output || '',
                error: error || null
            };
        } catch (error) {
            console.error("Python execution error:", error);
            return {
                output: '',
                error: error.message || 'An error occurred during execution'
            };
        }
    },

    async installPackage(packageName: string): Promise<void> {
        if (!pyodide) {
            await this.initialize();
        }

        try {
            await pyodide.loadPackage(packageName);
            console.log(`Package ${packageName} installed successfully`);
        } catch (error) {
            console.error(`Failed to install package ${packageName}:`, error);
            throw error;
        }
    }
};