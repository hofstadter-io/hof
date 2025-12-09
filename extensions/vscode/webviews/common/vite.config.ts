import { resolve } from 'path'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'
import react from '@vitejs/plugin-react-swc'
import dts from 'vite-plugin-dts'
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    tsconfigPaths(),
    dts({
      include: ['src'],
      entryRoot: resolve(__dirname, "src"),
      tsconfigPath: resolve(__dirname, "tsconfig.app.json"),
      rollupTypes: true,
   })
  ],
  resolve: { alias: { "@/": resolve(__dirname, "src/") } },
  build: {
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      formats: ['es', 'umd'],
      name: "authr-react-tanstack",
      fileName: (format) => `index.${format}.js`,
    },
    rollupOptions: {
      external: [
        "react", "react-dom", "tailwindcss",
        "@tanstack/react-query",
        "@tanstack/react-router",
        "@tanstack/react-form",
        "react-cookie",
        "radix-ui",
      ],
      output: {
        globals: {
          react: "React",
          "react-dom": "ReactDOM",
          tailwindcss: "tailwindcss",
        },
      },
    },
    sourcemap: true,
    emptyOutDir: true,
  },
})
