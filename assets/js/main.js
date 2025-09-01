import Alpine from 'alpinejs';
import persist from '@alpinejs/persist';
import flexImages from 'javascript-flex-images';
import lazySizes from 'lazysizes';
import 'lazysizes/plugins/parent-fit/ls.parent-fit';
import media from './components/media.js';
import scan from './components/scan.js';
import login from './components/login.js';
import folderMenu from './components/menu.js';
import './components/timeline.js';
import './components/map.js';
import './components/toasts.js';
import './components/notify.js';

window.lazySizes = window.lazySizes || {};

window.Alpine = Alpine;

window.numberWithCommas = function (x) {
  return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

window.plural = function (x, singular, plural) {
  if (x === 1) {
    return singular;
  } else {
    return plural;
  }
};
Alpine.plugin(persist);

document.addEventListener('alpine:init', () => {
  Alpine.store('fullscreen', false);
  Alpine.store('zoom', false);
  Alpine.store('data', {
    model: {
      camera: '',
      lens: '',
      term: '',
      type: '',
      rating: parseInt(0),
      folder: '',
      tag: '',
      software: '',
      focalLength35: parseFloat(0),
      // document these, should correspond to Filter type
      term_result: '',
      direction: 'desc',
      mediatype: '',
    },
    filter_ui: Alpine.$persist({
      open: false,
    }),
    memories_ui: Alpine.$persist({
      open: true,
      date: '',
    }),
  });
});

Alpine.data('media', media);
Alpine.data('scan', scan);
Alpine.data('login', login);
Alpine.data('menu', folderMenu);
Alpine.start();

document.querySelectorAll('.flex-images').forEach((flex) => {
  new flexImages({ selector: `#${flex.id}`, rowHeight: flex.dataset.size });
});

if (window.location.pathname.startsWith('/favorites')) {
  document.addEventListener('lazybeforeunveil', () => {
    lazySizes.autoSizer.checkElems();
  });
}
