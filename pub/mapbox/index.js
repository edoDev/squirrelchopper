//var gwcUrl = '/geowebcache/service/gmaps?layers=planet&zoom={z}&x={x}&y={y}&FORMAT=application/x-protobuf;type=mapbox-vector';
var gwcUrl = '/tile/{z}/{x}/{y}'
var baseUrl = '/demo/';
var styleBaseUrl = baseUrl + 'styles/osm-bright-gl-style/';
var spriteUrl = styleBaseUrl + 'sprite';
var glyphsUrl = baseUrl + 'fonts/{fontstack}/{range}.pbf';

// We are offline so we don't need a valid access token. But we do need a value in there.
mapboxgl.accessToken = '<your access token here>';

var map;

fetch(styleBaseUrl + 'style.json').then(function(response) {
  response.json().then(function(glStyle) {
    // MapBox js doesn't like relative paths. So rewrite the paths in the style to use and fully qualified URL.
    glStyle.sources.openmaptiles.tiles = [window.location.origin + gwcUrl];

    glStyle.sprite = window.location.origin + spriteUrl;
    glStyle.glyphs = window.location.origin + glyphsUrl;

    map = new mapboxgl.Map({
      container: 'map',
      style: '../styles/bright-v9/style.json'
    //  style: 'https://openmaptiles.github.io/osm-bright-gl-style/style-cdn.json'
//glStyle
    //           zoom: 13,
    //           center: [-122.447303, 37.753574]
    });
    map.addControl(new mapboxgl.NavigationControl());
  });
});
