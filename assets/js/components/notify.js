document.addEventListener('alpine:init', () => {
  Alpine.store('poll', {
    userName: '',
    status: false,
    async init() {
      let response = await fetch('/status');
      if (response.status === 200) {
        let data = await response.json();
        this.status = data;
      }

      this.poll();
    },
    async poll() {
      let backoff = 1000;
      while (true) {
        try {
          let response = await fetch('/poll');
          if (response.redirected) {
            break;
          } else if (response.status === 200) {
            let data = await response.json();
            if (data.message !== 'none') {
              Alpine.store('toasts').createToast(data.message, 'notice');
            }
            if (data.status === 'scanning') {
              this.status = true;
            } else if (data.status === 'complete') {
              this.status = false;
            }
          }
        } catch (error) {
          console.error('Polling error:', error);
          await new Promise((resolve) => setTimeout(resolve, backoff));
          backoff = Math.min(backoff * 2, 30000);
        }
      }
    },
  });
});
