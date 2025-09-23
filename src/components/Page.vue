<template>
  <div class="page-container">
    <section class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">
          <span>Hello, I'm</span><br>
          <span class="name-highlight">chrom</span>
        </h1>
        <p class="hero-subtitle">
          A software developer specializing in web development and DevOps.
        </p>
        <div class="social-links">
          <a href="https://github.com/chrom9103" target="_blank"class="social-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-github "><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"></path><path d="M9 18c-4.51 2-5-2-7-2"></path></svg>
          </a>
        </div>
        <div class="cta-buttons">
          <RouterLink to="/works" class="cta-button primary">View Works</RouterLink>
          <!-- <a href="/contact" class="cta-button secondary">Contact Me</a> -->
        </div>
      </div>
      <div class="hero-image">
        <img src="../assets/IMG_0051.jpg" alt="chrom profile" class="profile-image">
      </div>
    </section>

    <section class="about-section">
      <div class="about-content">
        <h2 class="section-heading">About Me</h2>
        <p class="about-text">
          青山学院大学 学士2年(在学中)。
        </p>
        <p class="about-text">
          主にフロントエンドやインフラに興味を持ち、幅広い技術スタックを学んでいます。
        </p>
      </div>
    </section>

    <section class="skills-section">
      <h2 class="section-heading">Skills</h2>
      <div class="skills-grid">
        <div class="skill-card">
          <h3 class="skill-title">Frontend</h3>
          <p class="skill-list">Vue.js, TypeScript, HTML, CSS</p>
        </div>
        <div class="skill-card">
          <h3 class="skill-title">Backend</h3>
          <p class="skill-list">Node.js, Python</p>
        </div>
        <div class="skill-card">
          <h3 class="skill-title">DevOps / Infra</h3>
          <p class="skill-list">Kubernetes, Docker, MicroK8s, Linux</p>
        </div>
      </div>
    </section>

    <section class="contact-section">
      <div class="contact-content">
        <h2 class="section-heading">Contact</h2>
        <p class="contact-text">
          下記のメールアドレスにお気軽にご連絡ください
        </p>
        <a href="mailto:your.email@example.com" class="cta-button primary">chrom@chrom.jp</a>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* Hero Section */
.hero-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3rem;
  padding: 4rem 0;
  border-bottom: 1px solid var(--color-border-soft);
}

@media (min-width: 768px) {
  .hero-section {
    flex-direction: row;
    justify-content: space-between;
    padding: 6rem 0;
  }
}

.hero-content {
  text-align: center;
  flex-basis: 60%;
}

@media (min-width: 768px) {
  .hero-content {
    text-align: left;
  }
}

.hero-title {
  font-size: 3rem;
  font-weight: 800;
  color: var(--color-heading);
  line-height: 1.2;
  margin-bottom: 1rem;
  animation: slide-in-top 1s ease-out;
}

.name-highlight {
  color: #6a0dad;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--color-text-mute);
  margin-bottom: 2rem;
  animation: fade-in 2s ease-out;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
  animation: slide-in-bottom 1.5s ease-out;
}

@media (min-width: 768px) {
  .social-links {
    justify-content: flex-start;
  }
}

.social-icon {
  font-size: 1.5rem;
  color: var(--color-text);
  transition: transform 0.3s ease;
}

.social-icon:hover {
  transform: translateY(-5px);
  color: #6a0dad;
}

.cta-buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  animation: fade-in 2.5s ease-out;
}

@media (min-width: 768px) {
  .cta-buttons {
    flex-direction: row;
  }
}

.cta-button {
  display: inline-block;
  padding: 0.8rem 8rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: bold;
  transition: all 0.3s ease;
  font-size: 1rem;
  text-align: center;
}

.primary {
  background-color: #f3f3f3;
  color: rgb(0, 0, 0);
  border: 2px solid transparent;
}

.primary:hover {
  background-color: #bebebe;
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.secondary {
  background-color: transparent;
  color: var(--color-text);
  border: 2px solid var(--color-border);
}

.secondary:hover {
  background-color: var(--color-border-soft);
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.hero-image {
  flex-basis: 40%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.profile-image {
  width: 250px;
  height: 250px;
  object-fit: cover;
  border-radius: 50%;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
  animation: scale-up 1.5s ease-out;
}

/* General Section Styles */
section {
  padding: 4rem 0;
  border-bottom: 1px solid var(--color-border-soft);
}

.section-heading {
  text-align: center;
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 2rem;
}

/* About Section */
.about-section {
  text-align: center;
}

.about-content {
  max-width: 800px;
  margin: 0 auto;
}

.about-text {
  font-size: 1.1rem;
  color: var(--color-text);
  line-height: 1.8;
}

/* Skills Section */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
  max-width: 1000px;
  margin: 0 auto;
}

.skill-card {
  background-color: var(--color-background-soft);
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease;
}

.skill-card:hover {
  transform: translateY(-5px);
}

.skill-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #6a0dad;
  margin-bottom: 1rem;
}

.skill-list {
  font-size: 1rem;
  color: var(--color-text-mute);
  line-height: 1.6;
}

/* Contact Section */
.contact-content {
  text-align: center;
  max-width: 600px;
  margin: 0 auto;
}

.contact-text {
  font-size: 1.1rem;
  color: var(--color-text);
  line-height: 1.8;
  margin-bottom: 2rem;
}

/* Animations */
@keyframes slide-in-top {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-in-bottom {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes scale-up {
  from {
    opacity: 0;
    transform: scale(0.8);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>