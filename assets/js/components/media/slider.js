import Glide from '@glidejs/glide';
import Panzoom from '@panzoom/panzoom';

// set initial slides on load
const initSlider = function (position) {
  this.slides = [];
  this.mediaItems.forEach((item, index) => {
    this.slides.push({
      id: index - 1,
      media: item,
    });
  });

  if (position) {
    this.slidePosition = position;
  } else {
    this.slidePosition = this.slides
      .map(function (x) {
        return x.media.hash;
      })
      .indexOf(this.media.hash);
  }

  this.renderSlides();
  this.handleClicks();

  this.initGlide(this.slidePosition);
};

const initGlide = function (startAt) {
  this.glide = new Glide('.glide', {
    startAt: startAt,
    perView: 1,
    rewind: false,
    // dragThreshold: 60,
    // swipeThreshold: 60,
  });

  this.glide.on('run.after', (e) => {
    if (e.direction === '>') {
      this.navigateNext();
    } else if (e.direction === '<') {
      this.navigatePrev();
    }
  });

  this.glide.mount();
};

// sync slides to state
const renderSlides = function () {
  let html = '';

  this.slides.forEach((s, i) => {
    let orientation = 'horizontal';
    if (s.media.ratio > 1) {
      orientation = 'vertical';
    }
    html += `<div class="glide__slide" data-hash="${s.media.hash}">
      <figure
        class="image-fit__figure image-fit__figure-${s.media.type}"
        itemprop="associatedMedia"
        itemscope=""
        itemtype="http://schema.org/ImageObject"
      >
        <div
          class="image-fit__container"
        >
        <template x-if="'${s.media.type}' === 'image'">
          <template x-if="!loaded.includes(parseInt(${s.media.hash}))">
            <div class=image-fit__loading>
              <svg width="36" height="36" fill="#fff" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><style>.spinner_ajPY{transform-origin:center;animation:spinner_AtaB .75s infinite linear}@keyframes spinner_AtaB{100%{transform:rotate(360deg)}}</style><path d="M12,1A11,11,0,1,0,23,12,11,11,0,0,0,12,1Zm0,19a8,8,0,1,1,8-8A8,8,0,0,1,12,20Z" opacity=".25"/><path d="M10.14,1.16a11,11,0,0,0-9,8.92A1.59,1.59,0,0,0,2.46,12,1.52,1.52,0,0,0,4.11,10.7a8,8,0,0,1,6.66-6.61A1.42,1.42,0,0,0,12,2.69h0A1.57,1.57,0,0,0,10.14,1.16Z" class="spinner_ajPY"/></svg>
            </div>
          </template>
        </template>
          <template x-if="'${s.media.type}' === 'video'">
            <video
              id="video-${s.media.hash}"
              class="image-fit__video image-fit__video-${orientation}"
              poster="/img/${s.media.hash}/${this.getMaxWidth(s.media)}"
              style="aspect-ratio: ${s.media.width}/${s.media.height}"
              controls
            >
            </video>
          </template>
          <template x-if="'${s.media.type}' === 'image'">
            <img
              class="image-fit__img"
              srcset="${s.media.srcset}"
              sizes="100vw"
              src="data:image/gif;base64,R0lGODlhQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="
            />
          </template>

        </div>
      </figure>
    </div>
    `;
    this.$refs.slides.innerHTML = html;

    if (s.media.type === 'video') {
      this.$nextTick(() => {
        // i === 0 when there are no next of prev items
        if (this.glide.index === i || i === 0) {
          const video = document.getElementById(`video-${s.media.hash}`);
          const hls = new Hls({
            debug: false,
            maxBufferLength: 3,
            // autoStartLoad: false,
          });
          hls.loadSource(`/transcode/${s.media.hash}/index.m3u8`);
          hls.attachMedia(video);
          hls.on(Hls.Events.MEDIA_ATTACHED, function () {
            // hls.stopLoad();
            // video.muted = true;
            // video.play();
          });
        }
      });
    }
  });
};

const updateSlides = function () {
  this.glide.destroy();
  this.slides = [];
  this.mediaItems.forEach((item, index) => {
    this.slides.push({
      id: index - 1,
      media: item,
    });
  });

  this.slidePosition = this.slides
    .map(function (x) {
      return x.media.hash;
    })
    .indexOf(this.media.hash);

  this.renderSlides();
  this.handleClicks();
  this.initGlide(this.slidePosition);
};

const handleClicks = function () {
  this.$nextTick(() => {
    if (this.gliding || this.media.type === 'video') return;
    const active = document.querySelector('.glide__slide--active .image-fit__img');
    let eventClickPending = 0;

    const handleClick = (e) => {
      // zoom out immediately on double click
      if (this.$store.zoom && !this.panzooming) {
        this.resetViewer(active);
      }

      // handle double click after waiting
      if (e.detail == 2 && eventClickPending != 0 && active && !this.gliding && !this.panzooming) {
        this.glide.disable();

        if (this.$store.fullscreen) {
          // zoom out
          if (this.$store.zoom) {
            this.resetViewer(active);
          } else {
            this.$store.zoom = !this.$store.zoom;
          }
        }

        // zoom in
        if (this.$store.zoom) {
          this.glide.disable();
          this.initViewer(active);
        }
      } else if (e.detail === 1 && eventClickPending == 0 && active && !this.gliding && !this.panzooming) {
        // handle single click after waiting
        if (!this.$store.fullscreen) {
          this.$store.fullscreen = true;
          clearTimeout(eventClickPending);
        } else {
          eventClickPending = setTimeout(() => {
            eventClickPending = 0;

            // toggle full screen
            if (!this.$store.zoom) {
              this.$store.fullscreen = !this.$store.fullscreen;
            } else {
              if (!this.panzooming) {
                this.resetViewer(active);
              }
            }
          }, 175);
          this.glide.enable();
        }
      }

      if (this.$store.fullscreen || this.$store.zoom) {
        window.scrollTo(0, 0);
      }
    };
    active.addEventListener('click', handleClick);
  });
};

const initViewer = function (active) {
  this.panzoom = Panzoom(active, {
    animate: true,
    // contain: 'outside', // not working, odd sizing on close
    maxScale: 10,
  });

  const largestImg = getLastSrcsetEntry(active);

  active.src = largestImg.url;
  active.srcset = '';

  this.panzoom.zoom(calculateScaleTo100Percent(largestImg.width, active.width));

  active.addEventListener('panzoomchange', () => {
    this.panzooming = true;
  });
  active.addEventListener('panzoompan', () => {
    this.panzooming = true;
  });
  active.addEventListener('panzoomend', () => {
    setTimeout(() => {
      this.panzooming = false;
    }, 50);
  });
  active.addEventListener('panzoomreset', () => {
    setTimeout(() => {
      this.panzooming = false;
    }, 50);
  });
};

const resetViewer = function (active) {
  this.$store.zoom = false;
  if (active) {
    active.srcset = this.media.srcset;
    active.src = '';
  }
  if (this.panzoom && this.panzoom.reset) {
    this.panzoom.reset();
    this.panzoom.resetStyle();
    this.panzoom.destroy();
  }
  this.glide.enable();
};

const calculateScaleTo100Percent = (thumbnailWidth, renderedWidth) => {
  return thumbnailWidth / renderedWidth;
};

const getLastSrcsetEntry = (imageElement) => {
  if (!imageElement || !imageElement.srcset) {
    throw new Error('Invalid image element or missing srcset attribute');
  }

  // Split the srcset into an array of entries
  const srcsetEntries = imageElement.srcset.split(',');
  // Get the last entry
  const lastEntry = srcsetEntries[srcsetEntries.length - 1].trim();

  // Find the last space in the entry
  const lastSpaceIndex = lastEntry.lastIndexOf(' ');
  // Split at the last space
  const url = lastEntry.substring(0, lastSpaceIndex).trim();
  const widthDescriptor = lastEntry.substring(lastSpaceIndex + 1).trim();

  const width = widthDescriptor ? parseInt(widthDescriptor.replace('w', ''), 10) : null;

  return { url, width };
};

export { initSlider, initGlide, renderSlides, updateSlides, handleClicks, resetViewer, initViewer };
