//var baseUrl = '../styles/bright-v9/';
var baseUrl = '../styles/osm-bright-gl-style/';

var tilegrid = ol.tilegrid.createXYZ({
  tileSize: 256,
  maxZoom: 16,
  loadZooms: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14]
});
var resolutions = tilegrid.getResolutions();
var layer = new ol.layer.VectorTile({
  source: new ol.source.VectorTile({
    attributions: '© <a href="https://www.mapbox.com/map-feedback/">Mapbox</a> ' +
      '© <a href="http://www.openstreetmap.org/copyright">' +
      'OpenStreetMap contributors</a>',
    format: new ol.format.MVT(),
    tileGrid: ol.tilegrid.createXYZ({maxZoom: 16, loadZooms: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14]}),
    tilePixelRatio: 16,
    url: '/tile/{z}/{x}/{y}.mvt'
    //url: '/geowebcache/service/gmaps?layers=planet&zoom={z}&x={x}&y={y}&FORMAT=application/x-protobuf;type=mapbox-vector'
    // url: 'http://54.172.109.225:8080/geowebcache/service/gmaps?layers=planet&zoom={z}&x={x}&y={y}&FORMAT=application/x-protobuf;type=mapbox-vector'
  })
  // renderMode: 'vector'
});
var map = new ol.Map({
  target: 'map',
  view: new ol.View({
    center: [0, 0],
    zoom: 2
  })
});

fetch(baseUrl + 'style.json').then(function(response) {
  response.json().then(function(glStyle) {
     olms.applyBackground(map, glStyle);
    glStyle.sprite = baseUrl + 'sprite';
     olms.applyStyle(layer, glStyle, 'openmaptiles').then(function() {
       map.addLayer(layer);
     });
//    var styleFunc = olms.getStyleFunction(glStyle, 'openmaptiles', resolutions);
//    layer.setStyle(styleFunc);
//    map.addLayer(layer);
  });
});
