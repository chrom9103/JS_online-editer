<script setup lang="ts">
import { RouterLink } from 'vue-router'

// Workクラスを定義
class Work {
  constructor(
    public id: number,
    public title: string,
    public category: string,
    public description: string,
    public technologies: string[],
    public features: string[],
    public image: string | null,
    public playLink: string | null,
    public gitLink: string
  ) {}
}

const works = [
   new Work(
    1,
    'JavaScript Playground',
    'Webアプリケーション',
    'A web-based code editor for JavaScript. This project demonstrates real-time code execution and a simple, intuitive user interface.',
    ['Vue.js', 'microk8s', 'Monaco Editor'],
    ['リアルタイムコード実行', '直感的なUI', 'ソース管理(開発中)'],
    '/works/スクリーンショット 2025-09-15 180030.png',
    '/works/playground',
    'https://github.com/chrom9103/JS_online-editer'
  ),
  new Work(
    2,
    'PiedPiper-bot',
    'Discord Bot (utility)',
    'A Python bot for discord. This project includes features such as collecting and analyzing Discord user activity logs and granting permissions.',
    ['python(discord.py)', 'microk8s', 'PVC, PV'],
    ['ログ収集・分析', 'DIscord権限付与(sudo権限)', '乱数生成 他utils'],
    null,
    null,
    'https://github.com/chrom9103/PiedPiper_bot'
  )
]
</script>

<template>
  <div class="works-page">
    <section class="works-section">
      <h2 class="section-title">My Projects</h2>
      <div class="works-grid">
        <div v-for="work in works" :key="work.id" class="project-card">
          <div class="project-card-link">
            <div class="card-image-wrapper">
              <img v-if="work.image" :src="work.image" :alt="work.title" class="card-image" />
              <div v-else class="no-image-placeholder">No Image</div>
              <div class="project-category-tag">{{ work.category }}</div>
            </div>
            <div class="card-content-wrapper">
              <h3 class="project-title">{{ work.title }}</h3>
              <p class="project-description">{{ work.description }}</p>

              <div class="project-meta">
              </div>

              <div class="project-technologies">
                <span v-for="(tech, i) in work.technologies" :key="i" class="tech-tag">
                  {{ tech }}
                </span>
              </div>

              <div class="project-features-list">
                <h4>主な機能:</h4>
                <ul>
                  <li v-for="(feature, i) in work.features" :key="i">
                    {{ feature }}
                  </li>
                </ul>
              </div>

              <div class="card-actions">
                <a v-if="work.playLink":href="work.playLink" target="_blank" rel="noopener noreferrer" class="play-button">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-external-link mr-2"><path d="M15 3h6v6"></path><path d="M10 14 21 3"></path><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path></svg>
                  Play
                </a>
                <div v-else class="play-button-null">null</div>
                <a :href="work.gitLink" target="_blank" rel="noopener noreferrer" class="git-button">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-github "><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"></path><path d="M9 18c-4.51 2-5-2-7-2"></path></svg>
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.works-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.section-title {
  text-align: center;
  font-size: 2.5rem;
  color: var(--color-heading);
  margin-bottom: 3rem;
  font-weight: 700;
}

.works-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 2.5rem;
}

.project-card {
  background-color: var(--color-background-soft);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: flex;
  flex-direction: column;
}

.project-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
}

.project-card-link {
  text-decoration: none;
  color: inherit;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.card-image-wrapper {
  position: relative;
  width: 100%;
  padding-bottom: 60%;
  background-color: var(--color-background-mute);
}

.card-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image-placeholder {
  position: absolute;
  top: 40%;
  left: 40%;
  font-size: 1.5rem;
  color: var(--color-text-mute);
}

.project-category-tag {
  position: absolute;
  top: 1rem;
  left: 1rem;
  background-color: #6a0dad;
  color: white;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: bold;
  z-index: 10;
}

.card-content-wrapper {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.project-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 0.5rem;
  line-height: 1.3;
}

.project-description {
  font-size: 0.95rem;
  color: var(--color-text-mute);
  margin-bottom: 1.5rem;
  line-height: 1.6;
  flex-grow: 1;
}

.project-meta {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
  color: var(--color-text);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.meta-item .icon {
  font-size: 1.1em;
}

.project-technologies {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.tech-tag {
  background-color: var(--color-border);
  color: var(--color-text);
  padding: 0.4rem 0.8rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 500;
}

.project-features-list {
  margin-top: 0.5rem;
  margin-bottom: 2rem;
}

.project-features-list h4 {
  font-size: 1rem;
  color: var(--color-heading);
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.project-features-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.project-features-list li {
  font-size: 0.95rem;
  color: var(--color-text-mute);
  margin-bottom: 0.3rem;
  position: relative;
  padding-left: 1.2rem;
}

.project-features-list li::before {
  content: '•';
  color: var(--color-accent);
  position: absolute;
  left: 0;
  font-weight: bold;
}


.card-actions {
  display: flex;
  gap: 1rem;
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border-soft);
}

.play-button {
  flex-grow: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #6a0dad;
  color: white;
  padding: 0.8rem 1.2rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: bold;
  font-size: 1rem;
  transition: background-color 0.3s ease;
  gap: 0.5rem;
}

.play-button:hover {
  background-color: #8a2be2;
}

.play-button-null {
  flex-grow: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-border);
  color: white;
  padding: 0.8rem 1.2rem;
  border-radius: 8px;
  text-decoration: none;
  color: var(--color-text-mute);
  font-weight: bold;
  font-size: 1rem;
  transition: background-color 0.3s ease;
  gap: 0.5rem;
}

.git-button {
  background-color: var(--color-background-mute);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  padding: 0.8rem 1.2rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

.git-button:hover {
  background-color: var(--color-border);
  border-color: var(--color-border-hover);
}

/* レスポンシブデザインの調整 */
@media (max-width: 768px) {
  .works-grid {
    grid-template-columns: 1fr;
  }
  .section-title {
    font-size: 2rem;
    margin-bottom: 2rem;
  }
}
</style>