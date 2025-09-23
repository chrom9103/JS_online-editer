<template>
  <div class="app-container">
    <header class="header">
      <div class="header-content">
        <div class="logo-group">
          <RouterLink to="/" class="logo-link">
            <img src="../assets/IMG_0051.jpg" alt="chrom logo" class="logo" />
            <h1 class="site-title">{{ props.title }}</h1>
          </RouterLink>
        </div>

        <nav v-if="props.showNavigation" class="navigation">
          <RouterLink to="/works" class="nav-link">Works</RouterLink>
        </nav>
      </div>
    </header>

    <main class="main-content">
      <slot />
    </main>

    <footer class="footer">
      <div class="footer-content">
        <p>© chrom All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

interface Props {
  title?: string
  showNavigation?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'chrom',
  showNavigation: true
})
</script>

<style scoped>

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* ヘッダーコンテナのスタイル */
.header {
  background-color: var(--color-background-soft, #fff);
  border-bottom: 1px solid var(--color-border, #ccc);
  width: 100%;
  box-sizing: border-box;
  padding: 0;
  flex-shrink: 0;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 2rem;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  position: relative;
  flex-wrap: nowrap;
}

/* ロゴとサイトタイトル */
.logo-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
}

.logo {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 50%;
  flex-shrink: 0;
}

.site-title {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--color-heading, #333);
  margin: 0;
}

/* ナビゲーションメニューの配置 */
.navigation {
  display: flex;
  gap: 2rem;
  flex-shrink: 0;
  margin-left: auto;
}

.nav-link {
  padding: 0.5rem 0;
  text-decoration: none;
  color: var(--color-text, #333);
  transition: color 0.3s;
}

.nav-link:hover {
  color: var(--color-accent, #007bff);
}

.nav-link.router-link-active {
  font-weight: bold;
}

.main-content {
  padding: 2rem;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  flex-grow: 1;
}

.footer {
  background-color: var(--color-background-soft);
  border-top: 1px solid var(--color-border);
  width: 100%;
  box-sizing: border-box;
  padding: 0;
  flex-shrink: 0;
}

.footer-content {
  text-align: center;
  color: var(--color-text-mute);
  padding: 1rem 2rem;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

/* レスポンシブデザインの調整 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }
  
  .logo-group {
    order: -1;
  }
  
  .social-icons {
    order: 1;
  }

  .navigation {
    order: 2;
    margin-left: 0;
    gap: 1.5rem;
  }
}
</style>