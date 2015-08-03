'use strict';

var express = require('express');
var path = require('path');
var models = require('./models');
var carts = models.carts;
var users = models.users;
var usercart = models.usercart;
var userchallenge = models.userchallenge;

var app = express();

app.configure(function() {
  app.use(express.logger('dev'));
  app.use(express.json());
  app.use(express.urlencoded());
  app.use(express.methodOverride());
  app.use(express.static(path.join(__dirname, 'public')));
  app.use(app.router);
});

app.post('/carts/gettrips', carts.getTrips);
app.get('/users/:name', users.findByName);
app.post('/users/', users.create);
app.post('/usercart/:name/addcartdata', usercart.addCartData);
app.post('/usercart/:name/removecartdata', usercart.removeCartData);
app.get('/usercart/:name/findallcarts', usercart.findAllCarts)
app.post('/userchallenge/:name/addchallenge', userchallenge.addChallenge);
app.post('/userchallenge/:name/removechallenges', userchallenge.removeChallenges);
app.get('/userchallenge/:name/findallchallenges', userchallenge.findAllChallenges);
app.get('/userchallenge/:name/findChallengeByID', userchallenge.findChallengeByID);
// app.post('/carts', carts.addKitten);
// app.put('/carts/:id', carts.updateKitten);
// app.delete('/carts/:id', carts.deleteKitten);

var port = process.env.PORT || 3000;
app.listen(port, function() {
  console.log('Listening on port ' + port);
});


