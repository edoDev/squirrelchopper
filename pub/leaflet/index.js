//var url = '/geowebcache/service/gmaps?layers=planet&zoom={z}&x={x}&y={y}&FORMAT=application/x-protobuf;type=mapbox-vector';
var url ='/tile/{z}/{x}/{y}'
var b = L.vectorGrid.protobuf(url);

var map = L.map('map').setView([0, 0], 2);
map.addLayer(b);
