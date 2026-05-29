'use strict';

// Gestion des likes/dislikes — appel fetch vers /api/react, mise à jour du DOM sans reload.
document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.reaction-bar').forEach(bar => {
    const id   = bar.dataset.id;
    const type = bar.dataset.type;
    const likeBtn    = bar.querySelector('.like-btn');
    const dislikeBtn = bar.querySelector('.dislike-btn');

    if (!likeBtn || !dislikeBtn) return;
    if (likeBtn.disabled) return; // visiteur anonyme

    likeBtn.addEventListener('click', () => react(id, type, 1, likeBtn, dislikeBtn));
    dislikeBtn.addEventListener('click', () => react(id, type, -1, likeBtn, dislikeBtn));
  });
});

async function react(targetId, targetType, value, likeBtn, dislikeBtn) {
  // désactiver les boutons pendant la requête
  likeBtn.disabled = true;
  dislikeBtn.disabled = true;

  try {
    const res = await fetch('/api/react', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_id: targetId, target_type: targetType, value }),
    });

    if (!res.ok) {
      if (res.status === 401) {
        window.location.href = '/login';
        return;
      }
      throw new Error('Erreur réseau');
    }

    const data = await res.json();
    updateReactionUI(likeBtn, dislikeBtn, data);
  } catch (err) {
    console.error('Réaction échouée :', err);
  } finally {
    likeBtn.disabled = false;
    dislikeBtn.disabled = false;
  }
}

function updateReactionUI(likeBtn, dislikeBtn, data) {
  const likeCount    = likeBtn.querySelector('.react-count');
  const dislikeCount = dislikeBtn.querySelector('.react-count');

  // animer le bouton cliqué
  const activeBtn = data.user_reaction === 1 ? likeBtn : data.user_reaction === -1 ? dislikeBtn : null;
  [likeBtn, dislikeBtn].forEach(btn => btn.classList.remove('active', 'bumping'));

  if (activeBtn) {
    activeBtn.classList.add('active');
    activeBtn.classList.add('bumping');
    activeBtn.addEventListener('animationend', () => activeBtn.classList.remove('bumping'), { once: true });
  }

  // mettre à jour les aria-pressed
  likeBtn.setAttribute('aria-pressed', data.user_reaction === 1 ? 'true' : 'false');
  dislikeBtn.setAttribute('aria-pressed', data.user_reaction === -1 ? 'true' : 'false');

  // mettre à jour les compteurs avec une micro-animation
  animateCount(likeCount, data.likes);
  animateCount(dislikeCount, data.dislikes);
}

function animateCount(el, newVal) {
  if (!el) return;
  const old = parseInt(el.textContent, 10);
  if (old === newVal) return;

  el.style.transform = 'scale(1.4)';
  el.style.transition = 'transform .15s ease-out';
  el.textContent = newVal;

  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      el.style.transform = 'scale(1)';
    });
  });
}
