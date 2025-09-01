import leaflet from 'leaflet';
import leafletmarkercluster from 'leaflet.markercluster';

document.addEventListener('alpine:init', () => {
  Alpine.data('map', (tileServer) => {
    return {
      defaultMap: tileServer === '/tiles/{z}/{x}/{y}.png',
      init() {
        const params = new URLSearchParams(window.location.search);

        let lat = 0;
        let lng = 0;
        if (params.get('lat') && params.get('lng')) {
          lat = parseFloat(params.get('lat'));
          lng = parseFloat(params.get('lng'));
        } else {
        }

        let zoom = 3;
        if (params.get('zoom')) {
          zoom = parseInt(params.get('zoom'));
        }

        let maxZoom = 4;
        if (!this.defaultMap) {
          maxZoom = 19;
        }
        var tiles = L.tileLayer(tileServer, {
            maxZoom: maxZoom,
          }),
          latlng = L.latLng(lat, lng);

        var map = L.map('map', { center: latlng, zoom: zoom, layers: [tiles] });

        var markers = L.markerClusterGroup();

        var LeafIcon = L.Icon.extend({
          options: {},
        });

        var icon = new LeafIcon({
          iconUrl: '/static/marker-icon.png',
          iconRetinaUrl: '/static/marker-icon.png',
          shadowUrl: '/static/marker-shadow.png',
          shadowRetinaUrl: '/static/marker-shadow.png',
          iconAnchor: [12, 41],
          popupAnchor: [0, -41],
        });

        if (window.map_data === null) {
          return;
        } else {
          for (var i = 0; i < map_data.length; i++) {
            var a = map_data[i];
            var title = a[2];
            var marker = L.marker(new L.LatLng(a[0], a[1]), { title: title, icon: icon });
            marker
              .bindPopup(`<a href="/media/${a[2]}"><img src="/img/${a[2]}/400" width="400" /></a>`, {
                minWidth: 300,
              })
              .openPopup();
            markers.addLayer(marker);
          }
        }

        map.addLayer(markers);
        if (lat === 0 && lng === 0) {
          map.fitBounds(map_data);
        }

        map.on('moveend', () => {
          const url = new URL(window.location.href);
          url.searchParams.set('lat', map.getBounds().getCenter().lat);
          url.searchParams.set('lng', map.getBounds().getCenter().lng);
          url.searchParams.set('zoom', map.getZoom());
          history.pushState(null, document.title, url.href);
        });
      },
    };
  });
});
