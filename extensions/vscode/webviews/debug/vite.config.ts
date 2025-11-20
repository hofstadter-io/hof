import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(),
    react()
  ],
  
  // Use relative paths for all assets in the build
  base: '',
  
  build: {
    // 1. Output the built files to the extension's 'out' directory
    outDir: '../../extension/webviews/debug',

    // 2. Prevent inlining assets as data: URIs (breaks CSP)
    assetsInlineLimit: 0,
    
    // 3. Define a predictable file structure for our sidebar to load
    rollupOptions: {
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
      },
    },
  },
})