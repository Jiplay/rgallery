import justifiedLayout from 'justified-layout';

document.addEventListener('alpine:init', () => {
  Alpine.data('timeline', () => ({
    sectionStore: {},
    sectionStates: {}, // app state - list of all sections and their state
    timeline: {},
    onThisDay: [],
    config: {
      containerWidth: window.innerWidth <= 900 ? window.innerWidth * 0.925 : window.innerWidth * 0.8,
      targetRowHeight: 150,
      segmentsMargin: 0,
      segmentTitle: 20,
      sectionMargin: 10,
      windowWidth: window.innerWidth,
    },
    folders: [],
    tags: [],
    sectionObserver: '',
    total: 0,
    max: 0,
    current: '',
    loading: [],
    mobile: {
      timeline: {
        open: false,
      },
    },
    loaded: false,
    async init() {
      this.showMemories();
      let url = new URL(window.location.href);
      this.$store.data.model.rating = parseInt(url.searchParams.get('rating') || 0);
      this.$store.data.model.term = url.searchParams.get('term') || '';
      this.$store.data.model.term_result = url.searchParams.get('term') || '';
      this.$store.data.model.type = url.searchParams.get('type') || '';
      this.$store.data.model.direction = url.searchParams.get('direction') || 'desc';
      this.$store.data.model.tag = url.searchParams.get('subject') || '';
      this.$store.data.model.folder = url.searchParams.get('folder') || '';
      this.$store.data.model.camera = url.searchParams.get('camera') || '';
      this.$store.data.model.lens = url.searchParams.get('lens') || '';
      this.$store.data.model.software = url.searchParams.get('software') || '';
      this.$store.data.model.focalLength35 = parseFloat(url.searchParams.get('focalLength35') || 0);

      await this.fetch(false);

      addEventListener('resize', async () => {
        if (this.config.windowWidth !== window.innerWidth) {
          this.config.containerWidth = window.innerWidth <= 900 ? window.innerWidth * 0.925 : window.innerWidth * 0.8;
          await this.loadTimeline();
        }
      });
      url = new URL(window.location.href);

      // check for date via query param on load and scroll to it
      if (url.searchParams.has('date')) {
        const date = url.searchParams.get('date');
        const target = date.substring(-1, date.lastIndexOf('-'));

        const uglyLoadScroller = setInterval(() => {
          const month = document.getElementById(target);
          if (month) {
            this.scrollToDate(date);
            clearInterval(uglyLoadScroller);
          } else {
            console.error('month not found', target);
          }
        }, 400);
      }

      // back button
      window.addEventListener('popstate', () => {
        this.goBack();
      });

      // back button query param change
      window.onpopstate = function (e) {
        if (!url.searchParams.has('date')) {
          setTimeout(function () {
            window.scrollTo(0, 0);
          }, 2);
        } else {
          this.scrollToDate(url.searchParams.get('date'));
        }
      };

      this.$nextTick(() => {
        this.getCurrent();
      });

      document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
          this.reset();
        }
      });
    },
    updateURL(pushState) {
      const url = new URL(window.location.href);

      if (this.$store.data.model.term !== '') {
        url.searchParams.set('term', this.$store.data.model.term);
      } else {
        url.searchParams.delete('term');
      }

      if (parseInt(this.$store.data.model.rating) === 0) {
        url.searchParams.delete('rating');
      } else if (parseInt(this.$store.data.model.rating) !== 0) {
        url.searchParams.set('rating', parseInt(this.$store.data.model.rating));
      }

      if (this.$store.data.model.direction === 'asc') {
        url.searchParams.set('direction', this.$store.data.model.direction);
      } else if (this.$store.data.model.direction === 'desc') {
        url.searchParams.delete('direction');
      }

      if (this.$store.data.model.type === 'image' || this.$store.data.model.type === 'video') {
        url.searchParams.set('type', this.$store.data.model.type);
      } else if (this.$store.data.model.type === '') {
        url.searchParams.delete('type');
      }

      if (this.$store.data.model.folder !== '') {
        url.searchParams.set('folder', this.$store.data.model.folder);
      } else {
        url.searchParams.delete('folder');
      }

      if (this.$store.data.model.tag !== '') {
        url.searchParams.set('subject', this.$store.data.model.tag);
      } else {
        url.searchParams.delete('subject');
      }

      if (pushState) {
        history.pushState(null, document.title, url.toString());
      } else {
        history.replaceState(null, document.title, url.toString());
      }
    },
    goBack() {
      const url = new URL(window.location.href);

      // get params to fetch
      const term = url.searchParams.get('term');
      if (term) {
        this.$store.data.model.term = term;
      } else {
        this.$store.data.model.term = '';
      }

      this.fetch(false);
    },
    async fetch(pushState) {
      this.loaded = false;
      let params = {
        format: 'json',
      };

      const modelParams = {
        camera: this.$store.data.model.camera,
        lens: this.$store.data.model.lens,
        term: this.$store.data.model.term,
        type: this.$store.data.model.type,
        folder: this.$store.data.model.folder,
        tag: this.$store.data.model.tag,
        software: this.$store.data.model.software,
        focalLength35: this.$store.data.model.focalLength35,
      };

      for (const [key, value] of Object.entries(modelParams)) {
        if ((typeof value === 'string' && value !== '') || (typeof value === 'number' && value !== 0)) {
          params[key] = value;
        }
      }

      if (parseInt(this.$store.data.model.rating) >= 2) {
        params.rating = this.$store.data.model.rating;
      }

      if (this.$store.data.model.direction === 'asc') {
        params.direction = this.$store.data.model.direction;
      }

      fetch('/?' + new URLSearchParams(params))
        .then((response) => response.json())
        .then((response) => {
          this.loaded = true;
          if (response.segment && response.segment.length > 0) {
            this.buildTimeline(response.segment);
            this.sectionStore = response.segment;
            this.total = response.total;
            this.direction = response.direction;
            this.$store.data.model.rating = parseInt(response.filter.rating);
            this.$store.data.model.term = response.filter.term;
            this.$store.data.model.term_result = response.filter.term;
            this.$store.data.model.type = response.filter.type;
            this.$store.data.model.direction = response.direction;
            this.$store.data.model.tag = response.filter.subject;
            this.$store.data.model.folder = response.filter.folder;
            this.$store.data.model.camera = response.filter.camera;
            this.$store.data.model.lens = response.filter.lens;
            this.$store.data.model.software = response.filter.software;
            this.$store.data.model.focalLength35 = parseFloat(response.filter.focalLength35);
            this.$store.data.filter_ui.open = false;
            window.scrollTo(0, 0);
            this.loadTimeline();
          } else {
            this.timeline = {};

            this.sectionStore = [];
            this.total = 0;
            this.$store.data.model.term_result = this.$store.data.model.term;
            document.querySelector('#grid').innerHTML = '';
            this.$store.data.filter_ui.open = false;
            window.scrollTo(0, 0);
          }
          if (pushState) {
            this.updateURL(true);
          }
        })
        .catch((error) => {
          console.error(error);
          this.loaded = true;
          this.$store.toasts.createToast('There was an error reaching the server.', 'error');
        });
    },
    scrollToDate(target) {
      const month = target.substring(-1, target.lastIndexOf('-'));
      // scroll to month
      const month_target = document.getElementById(month);
      if (month_target) {
        month_target.scrollIntoView();
        const uglyScroller = setInterval(() => {
          const day = document.getElementById(target);
          if (day) {
            // scroll to day
            day.scrollIntoView({ behavior: 'smooth', block: 'start' });
            const url = new URL(window.location.href);
            url.searchParams.set('date', target);
            history.pushState(null, document.title, url.toString());
            clearInterval(uglyScroller);
          } else {
            console.error('day not found (scrollToDate)', target);
            clearInterval(uglyScroller);
            this.scrollToDate(target);
          }
        }, 10);
      } else {
        console.error('month not found (scrollToDate)', month);
      }
    },

    async loadTimeline() {
      await this.populateGrid(document.getElementById('grid'), this.sectionStore);
    },

    getSections() {
      return this.sectionStore.map((section) => {
        return {
          sectionId: section.sectionId,
          totalItems: section.totalItems,
          totalSegments: section.segments.length,
        };
      });
    },

    // one segment per day
    getSegments(sectionId) {
      return this.sectionStore.find((section) => section.sectionId == sectionId).segments;
    },

    async populateGrid(gridNode, sections) {
      let sectionsHtml = '';
      let prevSectionEnd = this.config.sectionMargin;
      for (const section of sections) {
        let v = JSON.parse(JSON.stringify(section));
        this.sectionStates[v.sectionId] = {
          ...v,
          lastUpdateTime: -1,
          height: this.estimateSectionHeight(v),
          top: prevSectionEnd,
        };

        sectionsHtml += this.getDetachedSectionHtml(this.sectionStates[v.sectionId]);
        prevSectionEnd += this.sectionStates[section.sectionId].height;
      }
      gridNode.innerHTML = sectionsHtml;

      this.sectionObserver = new IntersectionObserver(
        (entries, observer) => {
          entries.forEach((entry) => {
            const sectionDiv = entry.target;

            if (this.sectionStates) {
              this.sectionStates[sectionDiv.id].lastUpdateTime = entry.time;
            }

            if (entry.isIntersecting) {
              window.requestAnimationFrame(() => {
                if (this.sectionStates[sectionDiv.id].lastUpdateTime === entry.time) {
                  this.populateSection(sectionDiv, this.getSegments(sectionDiv.id));
                }
              });
            } else {
              window.requestAnimationFrame(() => {
                if (this.sectionStates[sectionDiv.id].lastUpdateTime === entry.time) {
                  this.sectionStates[sectionDiv.id].active = false;
                  this.detachSection(sectionDiv, entry.time);
                }
              });
            }
          });

          let uglyHeight = 0;
          document.querySelectorAll('#grid > .section').forEach((e) => {
            uglyHeight += e.getBoundingClientRect().height;
          });
          if (this.$store.data.memories_ui.open && !this.showReset()) {
            uglyHeight += 300;
          }

          const grid = document.getElementById('grid');
          const timeline = document.getElementById('timeline-wrap');
          const body = document.querySelector('body');
          if (timeline && grid) {
            grid.style.height = uglyHeight + 'px';
            timeline.style.height = uglyHeight + 'px';
            body.style.height = uglyHeight + 'px';
          }
        },
        {
          rootMargin: '200px 0px',
        }
      );
      gridNode.querySelectorAll('.section').forEach(this.sectionObserver.observe.bind(this.sectionObserver));
    },

    getDetachedSectionHtml(sectionState) {
      return `<div id="${sectionState.sectionId}" class="section" style="width: ${this.config.containerWidth}px; height: ${sectionState.height}px; top: ${sectionState.top}px; left: 0px;"></div>`;
    },

    estimateSectionHeight(section) {
      const unwrappedWidth = (3 / 2) * section.totalItems * this.config.targetRowHeight * (7 / 10);
      const rows = Math.ceil(unwrappedWidth / this.config.containerWidth);
      const height = rows * this.config.targetRowHeight + section.segments.length * this.config.segmentTitle;

      return height;
    },

    populateSection(sectionDiv, segments) {
      let sectionId = sectionDiv.id;
      let segmentsHtml = '';
      let prevSegmentEnd = this.config.segmentsMargin;
      for (const segment of segments) {
        const segmentInfo = this.getSegmentHtmlAndHeight(segment, prevSegmentEnd);
        segmentsHtml += segmentInfo.html;
        prevSegmentEnd += segmentInfo.height + this.config.segmentsMargin;
      }

      sectionDiv.innerHTML = segmentsHtml;
      const newSectionHeight = prevSegmentEnd;
      const oldSectionHeight = this.sectionStates[sectionId].height;

      const heightDelta = newSectionHeight - oldSectionHeight;
      if (heightDelta == 0) {
        return;
      }

      this.sectionStates[sectionId].height = newSectionHeight;
      sectionDiv.style.height = `${newSectionHeight}px`;

      Object.keys(this.sectionStates).forEach((sectionToAdjustId) => {
        if (this.$store.data.model.direction === 'desc') {
          if (sectionToAdjustId >= sectionId) {
            return;
          }
        } else if (this.$store.data.model.direction === 'asc') {
          if (sectionToAdjustId <= sectionId) {
            return;
          }
        }

        this.sectionStates[sectionToAdjustId].top += heightDelta;

        const sectionToAdjustDiv = document.getElementById(sectionToAdjustId);
        if (sectionToAdjustDiv) {
          sectionToAdjustDiv.style.top = `${this.sectionStates[sectionToAdjustId].top}px`;
        }
      });

      // fix scroll when adding sections
      if (window.scrollY > this.sectionStates[sectionId].top) {
        window.scrollBy(0, heightDelta);
      }

      this.getCurrent();
    },

    getSegmentHtmlAndHeight(segment, top) {
      let sizes_new = [];
      const sizes = segment.i;
      sizes.forEach((size) => {
        sizes_new.push({ width: size[0], height: size[1] });
      });
      let geometry = justifiedLayout(sizes_new, this.config);
      for (i in geometry.boxes) {
        geometry.boxes[i].media = segment.i[i];
      }
      let tiles = geometry.boxes.map(this.getTileHtml, this).join('\n');

      geometry.containerHeight = geometry.containerHeight + this.config.segmentTitle;
      const d = new Date(segment.s);

      return {
        html: `<div id="${segment.s}" class="segment" style="width: ${this.config.containerWidth}px; height: ${
          geometry.containerHeight
        }px; top: ${top}px; left: 0px;">
      <h5 class="segment-title">${d.toUTCString().slice(0, 16)}</h5>
      ${tiles}
    </div>`,
        height: geometry.containerHeight,
      };
    },

    getTileHtml(box) {
      let tileClass = '';
      let backgroundCss = '';
      if (box.media[4]) {
        tileClass = ` tile-${box.media[4]}`;
      }
      if (box.media[3]) {
        backgroundCss = `background: ${box.media[3]}; `;
      }
      if (box.media[4] === 'v') {
        return `<div class="tile${tileClass}" style="${backgroundCss}width: ${box.width}px; height: ${
          box.height
        }px; left: ${box.left + 6}px; top: ${box.top + this.config.sectionMargin}px;"
        >
          <a href="/media/${box.media[2]}${window.location.search}">
            <img x-show="loading.includes(parseInt(${box.media[2]}))" class="timeline-loading" src="/static/loading.svg" alt="loading..." width="24" height="24">
            <video
              id="video-${box.media[2]}"
              data-hash="${box.media[2]}"
              class="image-fit__video"
              poster="/img/${box.media[2]}/800"
              style="aspect-ratio: ${box.width}/${box.height}"
              @mouseenter="hoverVideo"
              @mouseleave="hoverStop"
              autoplay="" muted="" playsinline=""
            >
          </a>
        </div>`;
      } else {
        return `<div class="tile${tileClass}" style="${backgroundCss}width: ${box.width}px; height: ${
          box.height
        }px; left: ${box.left + 6}px; top: ${box.top + this.config.sectionMargin}px;">
          <a href="/media/${box.media[2]}${window.location.search}">
            <img data-srcset="${this.generateSrcset(box.media[2], box.media[0])}" data-sizes="auto" class="lazyload" width="${box.width}" height="${
              box.height
            }" />
          </a>
    </div>`;
      }
    },

    generateSrcset(imageId, originalWidth) {
      const sizes = [200, 400, 800];

      // If original width is less than 800, include it as an option
      if (originalWidth < 800 && !sizes.includes(originalWidth)) {
        sizes.push(originalWidth);
      }

      return sizes
        .filter((size) => size <= originalWidth)
        .map((size) => `/img/${imageId}/${size} ${size}w`)
        .join(', ');
    },

    detachSection(sectionDiv) {
      sectionDiv.innerHTML = '';
      this.getCurrent();
    },
    buildTimeline(s) {
      this.timeline = {};
      s.forEach((r) => {
        const year = r.sectionId.split('-')[0];
        this.timeline[year] = [];
      });
      s.forEach((r) => {
        const year = r.sectionId.split('-')[0];
        if (this.max <= r.totalItems) {
          this.max = r.totalItems;
        }
        this.timeline[year].push(r);
      });
    },
    backspace() {
      if (this.$store.data.model.term === '') {
        this.fetch();
      }
    },

    showReset() {
      if (this.$store.data.model.term !== '') {
        return true;
      }
      return false;
    },
    setInitialParams() {
      if (this.$store.data.model.folder) {
        this.$store.data.model.folder = filter.folder;
        this.folders = [{ key: filter.folder }];
      }

      if (this.$store.data.model.subject) {
        this.$store.data.model.tag = this.$store.data.model.subject;
        this.tags = [{ key: this.$store.data.model.subject, value: this.$store.data.model.subject }];
      }
    },
    humanizeDate(date) {
      const d = new Date(date);
      const month = d.toLocaleString('default', { month: 'long' });
      return `${d.getFullYear()} ${month}`;
    },
    humanizeDateShort(date) {
      const d = new Date(date);
      const month = d.toLocaleString('default', { month: 'short' });
      return `${d.getFullYear()} ${month}`;
    },
    calcMax(sectionTotal) {
      return (sectionTotal / this.max) * 100;
    },
    compare(a, b) {
      if (a.distance < b.distance) {
        return -1;
      }
      if (a.distance > b.distance) {
        return 1;
      }
    },
    hideMemories() {
      const today = new Date().toISOString().slice(0, 10);
      this.$store.data.memories_ui.date = today;
      this.$store.data.memories_ui.open = false;
    },
    showMemories() {
      const today = new Date().toISOString().slice(0, 10);

      if (this.$store.data.memories_ui.date === today) {
        this.$store.data.memories_ui.open = false;
      } else {
        this.$store.data.memories_ui.open = true;
      }

      if (this.$store.data.memories_ui.open) {
        fetch('/onthisday')
          .then((response) => response.json())
          .then((response) => {
            this.onThisDay = response;
          })
          .catch((error) => {
            console.error(error);
            this.$store.toasts.createToast('There was an error reaching the server.', 'error');
          });
      }
    },
    getCurrent() {
      const sections = document.querySelectorAll('.segment');

      // https://jsfiddle.net/dperelman/510ws7c9/
      const visibleSections = new Map();
      let observer = new IntersectionObserver(
        (entries, observer) => {
          // Update new intersectionRatios.
          entries.forEach((entry) => {
            entry.target.querySelectorAll('.ratio').forEach((r) => (r.innerText = entry.intersectionRatio));
            if (entry.isIntersecting) {
              visibleSections.set(entry.target, entry.intersectionRatio);
            } else {
              visibleSections.delete(entry.target);
            }
          });

          let max = -1;
          let mostVisibleSection = null;
          for (const section of sections) {
            const intersectionRatio = visibleSections.get(section);
            if (intersectionRatio && intersectionRatio > max) {
              max = intersectionRatio;
              mostVisibleSection = section;
            }
          }
          if (mostVisibleSection) {
            this.current = mostVisibleSection.id.substring(0, 7);
          }
        },
        {
          threshold: [0, 0.25, 0.5, 0.75, 1],
        }
      );

      for (const section of sections) {
        observer.observe(section);
      }
    },
    hoverVideo(event) {
      const hash = parseInt(event.target.dataset.hash);
      this.loading.push(hash);
      const video = event.target;
      const hls = new Hls({
        debug: false,
        maxBufferLength: 3,
      });
      hls.loadSource(`/transcode/${hash}/index.m3u8`);
      hls.attachMedia(video);
      hls.on(window.Hls.Events.FRAG_BUFFERED, () => {
        this.loading = this.loading.filter((e) => e !== hash);
      });
    },
    hoverStop(event) {
      event.target.pause();
    },
    toggleType(type) {
      if (type === 'image' && this.$store.data.model.type === '') {
        this.$store.data.model.type = 'image';
      } else if (type === 'image' && this.$store.data.model.type === 'image') {
        this.$store.data.model.type = 'image';
      } else if (type === 'image' && this.$store.data.model.type === 'video') {
        this.$store.data.model.type = '';
      } else if (type === 'video' && this.$store.data.model.type === 'video') {
        this.$store.data.model.type = 'video';
      } else if (type === 'video' && this.$store.data.model.type === '') {
        this.$store.data.model.type = 'video';
      } else if (type === 'video' && this.$store.data.model.type === 'image') {
        this.$store.data.model.type = '';
      }
    },
    isFilterActive() {
      if (
        this.$store.data.model.camera !== '' ||
        this.$store.data.model.lens !== '' ||
        // this.$store.data.model.term !== '' ||
        this.$store.data.model.type !== '' ||
        this.$store.data.model.rating !== parseInt(0) ||
        this.$store.data.model.folder !== '' ||
        this.$store.data.model.tag !== '' ||
        this.$store.data.model.software !== '' ||
        this.$store.data.model.focalLength35 !== parseFloat(0) ||
        this.$store.data.model.term_result !== '' ||
        this.$store.data.model.direction !== 'desc' ||
        this.$store.data.model.mediatype !== ''
      ) {
        return true;
      }
      return false;
    },
    reset() {
      this.$store.data.model.camera = '';
      this.$store.data.model.lens = '';
      this.$store.data.model.term = '';
      this.$store.data.model.type = '';
      this.$store.data.model.rating = parseInt(0);
      this.$store.data.model.folder = '';
      this.$store.data.model.tag = '';
      this.$store.data.model.software = '';
      this.$store.data.model.focalLength35 = parseFloat(0);
      this.$store.data.model.term_result = '';
      this.$store.data.model.direction = 'desc';
      this.$store.data.model.mediatype = '';
      this.fetch(true);
    },
    showRating(rating) {
      this.$store.data.model.rating = parseInt(rating);
      this.fetch(true);
    },
  }));
});
