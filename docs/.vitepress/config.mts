import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Joker",
  description: "Joker documentation",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Documentation', link: '/docs/introduction' }
    ],

    sidebar: [
      {
        text: 'Overview',
        items: [
          { text: 'Introduction', link: '/docs/introduction'},
          { text: 'What this project is *NOT* about', link: '/docs/what-this-is-not'},
          { text: 'Available commands', link: '/docs/available-commands'},
        ],
      },
      {
        text: 'Jokerfile', link: '/docs/jokerfile',
        items: [
          { 
            text: 'Service',
            items: [
              { text: 'Overview', link: '/docs/jokerfile/service'}
            ]
          }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/toudi/joker' }
    ]
  }
})
