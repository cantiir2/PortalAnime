frontend/
├── public/                    # Static files
├── src/
│   ├── assets/               # Assets (images, styles)
│   │   └── main.css         # Global CSS
│   ├── components/          # Reusable components
│   │   └── VideoPlayer.vue  # Custom video player
│   ├── router/              # Vue Router configuration
│   │   └── index.js
│   ├── stores/              # Pinia stores
│   │   ├── auth.js         # Authentication store
│   │   └── content.js      # Content management store
│   ├── views/              # Page components
│   │   ├── Admin.vue
│   │   ├── Browse.vue
│   │   ├── Home.vue
│   │   ├── Login.vue
│   │   ├── Profile.vue
│   │   ├── Register.vue
│   │   └── Watch.vue
│   ├── App.vue             # Root component
│   └── main.js            # Application entry point
├── .env                   # Environment variables
├── index.html            # HTML entry point
├── package.json          # NPM dependencies
├── postcss.config.js     # PostCSS configuration
├── tailwind.config.js    # Tailwind CSS configuration
└── vite.config.js        # Vite configuration