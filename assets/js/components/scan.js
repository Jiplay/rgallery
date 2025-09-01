export default () => ({
  status: 'scan',
  scan(type) {
    let route = '';
    if (type === 'default') {
      route = 'scan';
    } else if (type === 'metadata') {
      route = 'metadatascan';
    } else if (type === 'deepscan') {
      route = 'deepscan';
    } else if (type === 'thumbscan') {
      route = 'thumbscan';
    }
    this.status = 'Scanning...';
    this.$store.poll.status = true;
    fetch(`/${route}`)
      .then((response) => {
        if (response.redirected) {
          this.$store.toasts.createToast('Scans can only be initiated by an admin user.', 'error');
          this.$store.poll.status = false;
          return;
        }

        return response.text();
      })
      .then((text) => {
        this.status = text;
      })
      .catch((error) => {
        this.$store.toasts.createToast('There was an error reaching the server.', 'error');
        console.error(error);
      });
  },
});
