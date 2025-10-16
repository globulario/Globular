// Globular landing page (external JS for CSP)
// Footer year
document.addEventListener('DOMContentLoaded', () => {
  const yearEl = document.getElementById('y');
  if (yearEl) yearEl.textContent = new Date().getFullYear();

  // Health ping
  const statusEl = document.getElementById('status');
  if (statusEl) {
    fetch('/healthz').then(r => {
      if (r.ok) statusEl.innerHTML = 'Status: <strong class="ok">healthy</strong>';
      else statusEl.innerHTML = 'Status: <strong class="bad">unhealthy</strong>';
    }).catch(() => {
      statusEl.innerHTML = 'Status: <strong class="bad">offline</strong>';
    });
  }
});
