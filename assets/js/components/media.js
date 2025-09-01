import { initMap } from './media/map.js';
import {
  initSlider,
  initGlide,
  renderSlides,
  updateSlides,
  handleClicks,
  resetViewer,
  initViewer,
} from './media/slider.js';

export default (media, previous, next, collection, slug, tileServer) => ({
  media: media,
  previous: previous,
  next: next,
  mediaItems: [...previous, media, ...next],
  collection: collection,
  slug: slug,
  scanning: 'Rescan',
  loading: false,
  map: {},
  glide: {},
  slidePosition: 1,
  slides: [],
  initSlider: initSlider,
  initGlide: initGlide,
  renderSlides: renderSlides,
  updateSlides: updateSlides,
  handleClicks: handleClicks,
  initMap: initMap,
  resetViewer: resetViewer,
  initViewer: initViewer,
  background_position: {
    x: 0,
    y: 0,
  },
  gliding: false,
  loaded: [],
  panzoom: {},
  panzooming: false,
  defaultMap: tileServer === '/tiles/{z}/{x}/{y}.png',
  init() {
    this.initSlider();

    this.$nextTick(() => {
      this.initMap(tileServer);
      this.detectLoaded();
      this.watchGlide();
    });

    this.$watch('media', () => {
      this.detectLoaded();
    });

    // back button
    window.onpopstate = () => {
      this.goBack();
    };

    document.addEventListener('keydown', (event) => {
      const active = document.querySelector('.glide__slide--active .image-fit__img');

      if (event.key === 'Escape') {
        this.$store.fullscreen = false;
        this.resetViewer(active);
      }

      if (event.key === 'z') {
        this.toggleZoom();
      }

      if (event.key === 'f') {
        this.toggleFullscreen();
      }
    });
  },
  navigateNext() {
    if (this.next && this.next[0]) {
      this.fetchImageData(this.media.hash, this.next[0].hash, true, 'right');
    }
    this.resetViewer();
  },
  navigatePrev() {
    if (this.previous && this.previous[this.previous.length - 1]) {
      this.fetchImageData(this.media.hash, this.previous[this.previous.length - 1].hash, true, 'left');
    }
    this.resetViewer();
  },
  navigateTo(oldHash, newImage, pushState, index) {
    this.fetchImageData(oldHash, newImage.hash, pushState, -1, index);
    this.resetViewer();
  },
  updateURL(oldHash, newHash, pushState) {
    let url = new URL(window.location.href);
    url = url.pathname.replace(oldHash, newHash);
    if (pushState) {
      history.pushState(null, document.title, url.toString() + window.location.search);
    } else {
      history.replaceState(null, document.title, url.toString() + window.location.search);
    }
  },
  fetchImageData(oldHash, newHash, pushState, dir, index) {
    if (this.loading === false) {
      this.loading = true;

      let filter = '';
      if (this.collection && this.slug) {
        filter = `/in/${collection}/${slug}`;
      } else if (this.collection) {
        filter = `/in/${collection}`;
      } else if (window.location.search) {
        filter = window.location.search;
      }

      fetch(`/media/${newHash}${filter}`, {
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then((response) => response.json())
        .then((data) => {
          this.media = data.media;
          this.previous = data.previous;
          this.next = data.next;
          this.mediaItems = [...data.previous, data.media, ...data.next];

          this.updateURL(oldHash, data.media.hash, pushState);
          document.title = data.media.path + ' | rgallery';
          this.loading = false;
          this.$nextTick(() => {
            this.initMap(tileServer);
          });

          if (this.$refs['video']) {
            this.$refs['video'].load();
          }

          this.updateSlides();
        })
        .catch((err) => {
          console.error(err);
          this.loading = false;
          this.$store.toasts.createToast('There was an error reaching the server.', 'error');
        });
    }
  },
  goBack() {
    let h = parseInt(window.location.pathname.split('/')[2]);
    if (h !== this.media.hash) {
      // find image by hash to navigate to
      const newImage = this.mediaItems.filter((img) => {
        return img.hash === h;
      });
      if (newImage && newImage.length) {
        this.navigateTo(h, newImage[0], false);
      } else {
        console.error('image not found');
      }
    } else {
      history.back();
    }
  },
  formatDay(date) {
    return date.split('T')[0];
  },
  async rescan() {
    this.scanning = 'Scanning...';
    const response = await fetch(`/media/${this.media.hash}`, { method: 'POST' });
    const data = await response.text();
    this.scanning = 'Scanned';
  },
  trimPrefix(prefix, s) {
    if (s.indexOf(prefix) === 0) {
      return s.substring(prefix.length);
    }
  },
  encodeTagKey(key) {
    return encodeURIComponent(key);
  },
  getMaxWidth(i) {
    const srcset = i.srcset.split(' ');
    return parseInt(srcset[srcset.length - 1].replace('w', ''));
  },
  // displayDate takes a local time and offset and returns a date string in RFC1123Z format.
  displayDate(d, o) {
    // https://stackoverflow.com/questions/7403486/add-or-subtract-timezone-difference-to-javascript-date
    const targetTime = new Date(d);
    if (o) {
      const timeZoneFromDB = o / 60; //time zone value from database
      //get the timezone offset from local time in minutes
      const tzDifference = timeZoneFromDB * 60 + targetTime.getTimezoneOffset();
      //convert the offset to milliseconds, add to targetTime, and make a new Date
      const offsetTime = new Date(targetTime.getTime() + tzDifference * 60 * 1000);

      const dayString = new Intl.DateTimeFormat('en-GB', {
        weekday: 'short',
        year: 'numeric',
        month: 'long',
        day: '2-digit',
      }).format(offsetTime);

      const timeString = new Intl.DateTimeFormat('en-GB', {
        hour: 'numeric',
        minute: 'numeric',
        second: 'numeric',
        fractionalSecondDigits: 3,
        hour12: false,
      }).format(offsetTime);

      return `${dayString} ${timeString}`;
    } else {
      // the UTC offset of the photo is unknown, so display it as is.
      const dayString = new Intl.DateTimeFormat('en-GB', {
        weekday: 'short',
        year: 'numeric',
        month: 'long',
        day: '2-digit',
        timeZone: 'GMT',
      }).format(targetTime);

      const timeString = new Intl.DateTimeFormat('en-GB', {
        hour: 'numeric',
        minute: 'numeric',
        second: 'numeric',
        fractionalSecondDigits: 3,
        hour12: false,
        timeZone: 'GMT',
      }).format(targetTime);

      return `${dayString} ${timeString}`;
    }
  },
  detectLoaded() {
    this.$nextTick(() => {
      document.querySelectorAll('.glide__slide').forEach((el) => {
        const image = el.querySelector('.image-fit__img');
        if (image) {
          if (!image.complete) {
            image.addEventListener('load', () => {
              this.loaded.push(parseInt(el.dataset.hash));
            });
            // image.addEventListener('error', handleImageLoad);
          } else {
            this.loaded.push(parseInt(el.dataset.hash));
          }
        }
      });
    });
  },
  watchGlide() {
    const attrObserver = new MutationObserver((mutations) => {
      mutations.forEach((mu) => {
        if (mu.type !== 'attributes' && mu.attributeName !== 'class') return;
        if (document.querySelector('.glide').classList.contains('glide--dragging')) {
          this.gliding = true;
        } else {
          setTimeout(() => {
            this.gliding = false;
          }, 100);
        }
      });
    });

    const gliders = document.querySelectorAll('.glide');
    gliders.forEach((el) => attrObserver.observe(el, { attributes: true }));
  },
  toggleFullscreen() {
    if (!this.panzooming) {
      this.$store.fullscreen = !this.$store.fullscreen;
      window.scrollTo(0, 0);
    }
  },
  toggleZoom() {
    const active = document.querySelector('.glide__slide--active .image-fit__img');

    if (this.media.type === 'video') return;
    if (!this.panzooming && !this.$store.zoom) {
      this.$store.fullscreen = true;
      this.$store.zoom = true;
      window.scrollTo(0, 0);

      // wait for fullscreen transition to finish before zooming
      setTimeout(() => {
        this.initViewer(active);
      }, 200);
    } else {
      this.resetViewer(active);
      this.$store.zoom = false;
      this.$store.fullscreen = false;
    }
  },
});
