import DefaultTheme from 'vitepress/theme'
import './tailwind.css'
import HomeLayout from './HomeLayout.vue'
import type { Theme } from 'vitepress'

export default {
  extends: DefaultTheme,
  Layout: HomeLayout
} satisfies Theme
