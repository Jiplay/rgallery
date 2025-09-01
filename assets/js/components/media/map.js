const initMap = function (tileServer) {
  const imageMetaContainer = document.querySelector('.image-meta__container');
  if (imageMetaContainer && imageMetaContainer.children.length >= 3) {
    imageMetaContainer.removeChild(imageMetaContainer.children[2]);
  }

  if (this.media.latitude && this.media.longitude) {
    let mapEl = `<div class="w-100">
              <div class="mb-1">
                <div class="image-meta__heading">Location</div>
                <a
                  class="image-meta__link"
                  href="/map?lat=${this.media.latitude}&lng=${this.media.longitude}&zoom=19"

                >${this.media.location}</a>
              </div>
              <div>
                <div id="map-container">
                  <div id="map" class="map" x-ref="map"></div>
                </div>
                ${
                  tileServer === '/tiles/{z}/{x}/{y}.png'
                    ? `<div class="body-small">
                  <p class="p-1">
                    For a higher resolution map, set the <code>--tile-server</code> CLI flag or the
                    <code>RGALLERY_TILE_SERVER</code> environmental variable with a tile server URL.
                  </p>
                </div>`
                    : ''
                }
              </div>
            </div>`;

    imageMetaContainer.insertAdjacentHTML('beforeend', mapEl);
  }

  let map = document.getElementById('map');

  if (map && this.media.latitude && this.media.longitude) {
    L.DomUtil.get('map');
    if (map != null) {
      map._leaflet_id = null;
    }

    this.map = L.map('map').setView([this.media.latitude, this.media.longitude], 13);

    let maxZoom = 4;
    if (!this.defaultMap) {
      maxZoom = 19;
    }
    L.tileLayer(tileServer, {
      maxZoom: maxZoom,
    }).addTo(this.map);

    let LeafIcon = L.Icon.extend({
      options: {},
    });

    const icon = new LeafIcon({
      iconUrl: '/static/marker-icon.png',
      iconRetinaUrl: '/static/marker-icon.png',
      shadowUrl: '/static/marker-shadow.png',
      shadowRetinaUrl: '/static/marker-shadow.png',
      iconAnchor: [12, 41],
      popupAnchor: [0, -41],
    });

    let marker = L.marker([this.media.latitude, this.media.longitude], { icon: icon }).addTo(this.map);
    marker.url = '/map';
    marker.on('click', function () {
      window.location = `/map?lat=${this.media.latitude}&lng=${this.media.longitude}&zoom=19`;
    });
  }
};

export { initMap };
