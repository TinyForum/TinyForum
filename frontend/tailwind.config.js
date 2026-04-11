/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-inter)', 'system-ui', 'sans-serif'],
        mono: ['var(--font-fira-code)', 'monospace'],
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        light: {
          'primary': '#ef4444',           // 红色主色
          'primary-focus': '#dc2626',     // 深红色（hover状态）
          'primary-content': '#ffffff',   // 主色上的文字颜色
          
          'secondary': '#f97316',         // 橙红色作为辅助色
          'secondary-focus': '#ea580c',
          'secondary-content': '#ffffff',
          
          'accent': '#ec4899',            // 粉红色作为强调色
          'accent-focus': '#db2777',
          'accent-content': '#ffffff',
          
          'neutral': '#4b5563',           // 中性灰
          'neutral-focus': '#374151',
          'neutral-content': '#ffffff',
          
          'base-100': '#ffffff',          // 基础背景色（浅色模式）
          'base-200': '#fef2f2',          // 浅红色背景
          'base-300': '#fee2e2',          // 更浅的红色背景
          'base-content': '#1f2937',      // 基础文字颜色
          
          'info': '#3b82f6',
          'success': '#10b981',
          'warning': '#f59e0b',
          'error': '#ef4444',
          
          // 自定义红色主题特有颜色
          'red-50': '#fef2f2',
          'red-100': '#fee2e2',
          'red-200': '#fecaca',
          'red-300': '#fca5a5',
          'red-400': '#f87171',
          'red-500': '#ef4444',
          'red-600': '#dc2626',
          'red-700': '#b91c1c',
          'red-800': '#991b1b',
          'red-900': '#7f1d1d',
        },
        dark: {
          'primary': '#f87171',           // 亮红色（深色模式）
          'primary-focus': '#ef4444',
          'primary-content': '#ffffff',
          
          'secondary': '#fb923c',         // 亮橙红色
          'secondary-focus': '#f97316',
          'secondary-content': '#ffffff',
          
          'accent': '#f472b6',            // 亮粉红色
          'accent-focus': '#ec4899',
          'accent-content': '#ffffff',
          
          'neutral': '#374151',
          'neutral-focus': '#1f2937',
          'neutral-content': '#f3f4f6',
          
          'base-100': '#1f2937',          // 深色背景
          'base-200': '#111827',          // 更深背景
          'base-300': '#0f172a',          // 最深背景
          'base-content': '#f3f4f6',      // 深色模式文字颜色
          
          'info': '#60a5fa',
          'success': '#34d399',
          'warning': '#fbbf24',
          'error': '#f87171',
          
          // 深色模式红色系
          'red-400': '#f87171',
          'red-500': '#ef4444',
          'red-600': '#dc2626',
          'red-700': '#b91c1c',
        },
      },
      // 可选：纯红色主题（更激进）
      {
        'red-dark': {
          'primary': '#dc2626',
          'primary-focus': '#b91c1c',
          'primary-content': '#ffffff',
          'secondary': '#991b1b',
          'secondary-focus': '#7f1d1d',
          'secondary-content': '#ffffff',
          'accent': '#ef4444',
          'accent-focus': '#dc2626',
          'accent-content': '#ffffff',
          'neutral': '#2d2d2d',
          'neutral-focus': '#1f1f1f',
          'neutral-content': '#ffffff',
          'base-100': '#1a1a1a',
          'base-200': '#2d2d2d',
          'base-300': '#404040',
          'base-content': '#e5e5e5',
          'info': '#3b82f6',
          'success': '#10b981',
          'warning': '#f59e0b',
          'error': '#ef4444',
        },
      },
    ],
    darkTheme: 'dark',  // 默认深色主题
    base: true,
    styled: true,
    utils: true,
    logs: false,
  },
};