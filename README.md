# ✨ PortalAnime

[![Vue.js](https://img.shields.io/badge/Vue-Green?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![Go](https://img.shields.io/badge/Go-Blue?style=flat-square&logo=go)](https://go.dev/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-3ca9d1?style=flat-square&logo=tailwind-css)](https://tailwindcss.com/)
[![Vite](https://img.shields.io/badge/Vite-646CFF?style=flat-square&logo=vitejs)](https://vitejs.dev/)
[![npm](https://img.shields.io/badge/npm-CB3837?style=flat-square&logo=npm)](https://www.npmjs.com/)


> PortalAnime adalah aplikasi web halaman tunggal (SPA) yang dibangun menggunakan Vue.js di frontend dan Go di backend.  Aplikasi ini menyediakan platform untuk menonton dan mengelola anime.

## ✨ Fitur Utama

* **Manajemen Konten Anime:**  Aplikasi ini memungkinkan pengelolaan berbagai aspek konten anime, termasuk detailnya, episode, musim, dan genre.
* **Sistem Otentikasi Pengguna:**  Pengguna dapat mendaftar, masuk, dan keluar dari platform dengan aman, memungkinkan pengalaman yang dipersonalisasi.
* **Penelusuran dan Pencarian Anime:**  Pengguna dapat menelusuri dan mencari konten anime berdasarkan genre, judul, dan kriteria lainnya.
* **Detail Anime dan Episode:**  Tampilan detail yang komprehensif untuk setiap anime dan episode, termasuk informasi terkait.
* **Riwayat Tontonan:**  Fitur untuk melacak dan mengelola riwayat tontonan pengguna.
* **Unggah Video:**  Kemungkinan untuk mengunggah video anime baru ke platform (berdasarkan keberadaan `VideoUploader.vue`).
* **Pemutar Video:** Pemutar video terintegrasi untuk memutar konten anime.
* **Administrasi:** Fitur administrasi yang mungkin memungkinkan kontrol atas seluruh konten dan pengguna (berdasarkan keberadaan `Admin.vue`).


## 🛠️ Tumpukan Teknologi

| Kategori          | Teknologi       | Catatan                                      |
|----------------------|-------------------|----------------------------------------------|
| Bahasa Frontend     | Vue.js           | Framework JavaScript untuk frontend           |
| Bahasa Backend      | Go               | Bahasa pemrograman untuk backend              |
| Build Tool Frontend | Vite             | Tool untuk membangun frontend                   |
| Styling Frontend   | Tailwind CSS     | Framework CSS untuk styling frontend           |
| Manajer Paket      | npm              | Manajer paket untuk dependensi frontend       |
| Database           | N/A              | Tipe database tidak ditentukan secara eksplisit |
| Routing (Frontend) | Vue Router       | Library routing untuk Vue.js                  |
| State Management (Frontend) | Pinia           | Library state management untuk Vue.js         |


## 🏛️ Tinjauan Arsitektur

Aplikasi ini mengikuti arsitektur client-server. Frontend (Vue.js) berkomunikasi dengan backend (Go) melalui API RESTful. Backend menangani logika bisnis, akses database, dan manajemen data.  Frontend bertanggung jawab atas presentasi dan interaksi pengguna.


## 🚀 Memulai

1. **Clone repositori:**
   ```bash
   git clone https://github.com/cantiir2/PortalAnime.git
   ```

2. **Navigasi ke direktori frontend:**
   ```bash
   cd PortalAnime/frontend
   ```

3. **Instal dependensi:**
   ```bash
   npm install
   ```

4. **Jalankan server pengembangan:**
   ```bash
   npm run dev
   ```

   Untuk menjalankan backend, langkah-langkah yang serupa diperlukan tetapi untuk direktori `backend` menggunakan `go mod`. Instruksi detailnya belum disediakan.

## 📂 Struktur File

```
/
├── LICENSE
├── README.md
├── backend
│   ├── .air.toml
│   ├── .env
│   ├── .env copy
│   ├── .gitignore
│   ├── cmd
│   │   └── server
│   │       ├── .env
│   │       └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── api
│   │   │   ├── handlers
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── content_handler.go
│   │   │   │   ├── episode_handler.go
│   │   │   │   ├── episode_handler_test.go
│   │   │   │   ├── genre_handler.go
│   │   │   │   ├── media_handler.go
│   │   │   │   ├── season_handler.go
│   │   │   │   └── watch_history_handler.go
│   │   │   ├── middleware
│   │   │   │   └── auth.go
│   │   │   └── routes
│   │   │       └── routes.go
│   │   ├── config
│   │   │   └── config.go
│   │   ├── db
│   │   │   └── db.go
│   │   ├── models
│   │   │   ├── category.go
│   │   │   ├── content.go
│   │   │   ├── content_type.go
│   │   │   ├── episode.go
│   │   │   ├── genre.go
│   │   │   ├── season.go
│   │   │   ├── user.go
│   │   │   └── watch_history.go
│   │   ├── repository
│   │   │   ├── category_repository.go
│   │   │   ├── content_repository.go
│   │   │   ├── episode_repository.go
│   │   │   ├── genre_repository.go
│   │   │   ├── season_repository.go
│   │   │   ├── user_repository.go
│   │   │   └── watch_history_repository.go
│   │   └── services
│   │       ├── content_service.go
│   │       ├── episode_service.go
│   │       ├── genre_service.go
│   │       ├── media_service.go
│   │       ├── season_service.go
│   │       ├── user_service.go
│   │       └── watch_history_service.go
│   └── struckture-beckend.md
└── frontend
    ├── .env
    ├── .gitignore
    ├── .vite
    │   └── deps
    │       ├── _metadata.json
    │       └── package.json
    ├── .vscode
    │   └── extensions.json
    ├── README.md
    ├── frontend@0.0.0
    ├── index.html
    ├── jsconfig.json
    ├── npm
    ├── package-lock.json
    ├── package.json
    ├── postcss.config.js
    ├── public
    │   └── favicon.ico
    ├── src
    │   ├── App.vue
    │   ├── assets
    │   │   ├── base.css
    │   │   ├── logo.svg
    │   │   └── main.css
    │   ├── components
    │   │   ├── HelloWorld.vue
    │   │   ├── TheWelcome.vue
    │   │   ├── VideoPlayer copy.vue
    │   │   ├── VideoPlayer.vue
    │   │   ├── VideoUploader.vue
    │   │   ├── WelcomeItem.vue
    │   │   └── icons
    │   │       ├── IconCommunity.vue
    │   │       ├── IconDocumentation.vue
    │   │       ├── IconEcosystem.vue
    │   │       ├── IconSupport.vue
    │   │       └── IconTooling.vue
    │   ├── main.js
    │   ├── router
    │   │   └── index.js
    │   ├── stores
    │   │   ├── auth.js
    │   │   └── content.js
    │   └── views
    │       ├── Admin.vue
    │       ├── Browse.vue
    │       ├── ContentDetail.vue
    │       ├── EpisodeDetail.vue
    │       ├── Home.vue
    │       ├── Login.vue
    │       ├── Musim.vue
    │       ├── Profile.vue
    │       ├── Register.vue
    │       └── Watch.vue
    ├── struckture-frontend.md
    ├── tailwind.config.js
    └── vite.config.js
```

* **backend:** Direktori ini berisi kode sumber untuk backend Go, termasuk handler API, model data, repositori, dan layanan.
* **frontend:** Direktori ini berisi kode sumber untuk frontend Vue.js, termasuk komponen, routing, dan store.


