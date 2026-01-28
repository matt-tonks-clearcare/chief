import { defineConfig } from 'vitepress'
import tailwindcss from '@tailwindcss/vite'
import fs from 'node:fs'
import path from 'node:path'

export default defineConfig({
  title: 'Chief',
  description: 'Autonomous PRD Agent',
  base: '/chief/',

  // Force dark mode only
  appearance: 'force-dark',

  vite: {
    plugins: [tailwindcss()]
  },

  markdown: {
    theme: 'tokyo-night'
  },

  async transformPageData(pageData, { siteConfig }) {
    const filePath = path.join(siteConfig.srcDir, pageData.relativePath)
    try {
      const rawContent = fs.readFileSync(filePath, 'utf-8')
      pageData.frontmatter.head ??= []
      pageData.frontmatter.head.push([
        'script',
        {},
        `window.__DOC_RAW = ${JSON.stringify(rawContent)};`
      ])
    } catch {
      // File not found — skip injection
    }
  },

  themeConfig: {
    siteTitle: 'Chief',

    nav: [
      { text: 'Home', link: '/' },
      { text: 'Docs', link: '/guide/quick-start' },
      { text: 'GitHub', link: 'https://github.com/minicodemonkey/chief' }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/minicodemonkey/chief' }
    ],

    search: {
      provider: 'local'
    },

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Quick Start', link: '/guide/quick-start' },
          { text: 'Installation', link: '/guide/installation' }
        ]
      },
      {
        text: 'Concepts',
        items: [
          { text: 'How Chief Works', link: '/concepts/how-it-works' },
          { text: 'The Ralph Loop', link: '/concepts/ralph-loop' },
          { text: 'PRD Format', link: '/concepts/prd-format' },
          { text: 'The .chief Directory', link: '/concepts/chief-directory' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'CLI Commands', link: '/reference/cli' },
          { text: 'Configuration', link: '/reference/configuration' },
          { text: 'PRD Schema', link: '/reference/prd-schema' }
        ]
      },
      {
        text: 'Troubleshooting',
        items: [
          { text: 'Common Issues', link: '/troubleshooting/common-issues' },
          { text: 'FAQ', link: '/troubleshooting/faq' }
        ]
      }
    ]
  }
})
