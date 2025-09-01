function menu() {
  return {
    toggle(e) {
      e.target.parentElement.parentElement.querySelector('.menu--sub').classList.toggle('hidden');
      e.target.parentElement.parentElement.classList.toggle('menu--open');
    },
  };
}

export default menu;
