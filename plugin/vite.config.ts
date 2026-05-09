import { defineConfig } from "vite";

export default defineConfig({
  build: {
    lib: {
      entry: "src/index.tsx",
      name: "MyPlugin",
      fileName: "index",
      formats: ["iife"], // 必须是 iife，才能在 <script> 标签中直接运行
    },
    rollupOptions: {
      // React 由宿主页面提供，不打包进来（减小体积）
      external: ["react", "react-dom"],
      output: {
        globals: {
          react: "React",
          "react-dom": "ReactDOM",
        },
      },
    },
  },
});
