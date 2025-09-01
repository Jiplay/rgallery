export default () => ({
  error: '',
  init() {
    const params = new URLSearchParams(window.location.search);
    this.error = params.get('error');
  },
});
