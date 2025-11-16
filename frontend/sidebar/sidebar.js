fetch('../sidebar/sidebar.html')
  .then(res => res.text())
  .then(html => {
    document.getElementById('sidebar-container').innerHTML = html;

    const hamburger = document.getElementById('hamburger');
    const sidebar = document.getElementById('sidebar');

    if (hamburger && sidebar) {
      hamburger.addEventListener('click', () => sidebar.classList.toggle('open'));
    }
  });
