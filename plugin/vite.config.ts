import { defineConfig } from 'vite';
import path from 'path';
import fs from 'fs';
import react from '@vitejs/plugin-react';
import copy from 'rollup-plugin-copy';

// 使用 process.cwd() 获取当前工作目录的路径
const manifestPath = path.resolve(process.cwd(), 'manifest.json');
const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf-8'));

export default defineConfig({
  plugins: [react()],
  build: {
    lib: {
      entry: path.resolve(process.cwd(), 'src/index.tsx'),
      name: manifest.id,
      formats: ['iife'],
      fileName: () => 'main.js',
    },
    rollupOptions: {
      external: manifest.external || [],
      output: {
        globals: manifest.globals || {},
      },
      plugins: [
        copy({
          targets: [
            { src: 'manifest.json', dest: 'dist' }
          ],
          hook: 'writeBundle'   // 在写完 bundle 后执行复制
        })
      ]
    },
    outDir: 'dist',
    emptyOutDir: true,
    
  },
   define: {
    '__PLUGIN_ID__': JSON.stringify(manifest.id)   ,// 假设 manifest.id = "14"
    'process.env': JSON.stringify({}),
    'process.env.NODE_ENV': JSON.stringify('production')
  }
});