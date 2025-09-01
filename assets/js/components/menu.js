function menu() {
  return {
    toggle(path) {
      document.querySelector(`[data-path="${path}"]`).classList.toggle('menu--open');
      document.querySelector(`[data-path="${path}"]`).querySelector('.menu--sub').classList.toggle('hidden');
    },
    toggleChild(path) {
      document.querySelector(`[data-path="${path}"]`).classList.toggle('hidden');
      document.querySelector(`[data-path="${path}"]`).parentElement.classList.toggle('menu--open');
    },
  };
}

export default menu;
