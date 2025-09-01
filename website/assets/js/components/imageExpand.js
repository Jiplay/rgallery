function imageExpand() {
  return {
    expanded: false,
    toggleImage: function () {
      const container = this.$root;
      const image = this.$root.querySelector('img');
      const button = this.$root.querySelector('.expand-btn');

      if (!this.expanded) {
        container.style.height = image.scrollHeight + 'px';
        button.classList.add('rotated');
      } else {
        container.style.height = '400px';
        button.classList.remove('rotated');
      }

      this.expanded = !this.expanded;
    },
  };
}

export default imageExpand;
