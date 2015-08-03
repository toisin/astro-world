'use strict';

var fs = require('fs');
var mongoose = require('mongoose');
var Schema = mongoose.Schema;

var mongoUri = process.env.MONGOLAB_URI ||
    process.env.MONGOHQ_URL ||
    'mongodb://localhost/cart_db';

mongoose.connect(mongoUri);
var db = mongoose.connection;


// cartSchema represents the collection of data that represents the relationships between
//     the Independent Variables and the Dependent Varaible
var cartSchema = mongoose.Schema({
  trips: [Number],
  handleLength: {type: String, enum: ['Long', 'Short'], required: true},
  wheelSize: {type: String, enum: ['Large(4)', 'Small(3)'], required: true},
  bucketSize: {type: String, enum: ['Big(13)', 'Small(10)'], required: true},
  bucketPlacement: {type: String, enum: ['Far', 'Near'], required: true},
});

var Cart = mongoose.model('Cart', cartSchema);


function getTrips(req, res) {
  console.log(req.body);
  Cart.find({handleLength:req.body.handleLength, bucketPlacement:req.body.bucketPlacement,
          bucketSize:req.body.bucketSize, wheelSize:req.body.wheelSize}, getTripsCallback(res));
  // Cart.findById(req.params.id, makeCallback(res));
}


function getTripsCallback(response) {
  return function(error, result) {
    var carts = result;
    var x = Math.floor((Math.random()*carts.length));
    if (error)
      response.send(500, {error: error});
    carts[x]
    response.send(200, filterHidden(carts[x]));
  };
}


exports.carts = {
  getTrips: getTrips
};


















// exports.findAll = function(req, res) {
//   Cart.find({}, function(err, result) {
//     res.send(result);
//   });
// };

// exports.findById = function(req, res) {
//   Cart.findById(req.params.id, makeCallback(res));
// };

// exports.addCart = function(req, res) {
//   var cart = new Cart(req.body);
//   cart.save(makeCallback(res));
// };

// exports.updateCart = function(req, res) {
//   Cart.findByIdAndUpdate(req.params.id, req.body, makeCallback(res));
// };

// exports.deleteCart = function(req, res) {
//   Cart.findByIdAndRemove(req.params.id, makeCallback(res));
// };


// exports.Cart = Cart;
function importDB(filename) {
  //check if db exists
  Cart.find({}, function (err, count){
    if (count<1) {
      var stream = fs.readFileSync(filename, 'utf-8');
      //console.log(stream);
      var lines = stream.split(/\n|\r/);

      for (var i = 1; i < lines.length; i++) {
        var t = lines[i].split(',');
        var c = new Cart({
          trips: [t[4]],
          handleLength: t[0],
          wheelSize: t[3],
          bucketSize: t[2],
          bucketPlacement: t[1],
        });
        c.save(function(err, result) {
          if(err)
           console.error(err);
          else
           console.log(result)
        });
      }      
    }
  });
};

// Import only the first time
importDB('casescart.csv');







// userCartSchema represents a cart that the user has saved to their notebook
var userCartSchema = mongoose.Schema({
  name: {type: String, required: true}, //username for the user this cart is saved for
  trips: [Number], //an array of number of trips that this cart could make
  handleLength: {type: String, enum: ['Long', 'Short'], required: true},
  wheelSize: {type: String, enum: ['Large(4)', 'Small(3)'], required: true},
  bucketSize: {type: String, enum: ['Big(13)', 'Small(10)'], required: true},
  bucketPlacement: {type: String, enum: ['Far', 'Near'], required: true},
  cartNumber: {type: Number, required: true}, //a unique number identifying this cart
                                              //    based on the selection made for the four IVs
                                              //    this number should be no larger than 16.
});

var UserCart = mongoose.model('UserCart', userCartSchema);

// define the function that will be called once the DB has looked
//     up a specific cart to see if the cart has previously been saved,
//     if it hasn't been, save a new one, otherwise, add an additional
//     trip to the existing cart.
function addTripCallback(response, newCart, totalCarts) {
  return function(error, cart) {
    if (cart) {
      cart.trips.push(newCart.trips[0]);
      cart.save(saveCallback(response));
    } else {
      newCart.cartNumber = totalCarts + 1;
      newCart.save(saveCallback(response));
    }
  };
}

// define the function that will be called once the DB returns 
//     the count of the number of carts this user has saved in order
//     to either add a new cart or add an additional trip to an existing
//     saved cart.
function addTripAfterCountCallBack(req, response, newCart) {
  return function(error, totalCarts) {
    console.log(totalCarts);

    UserCart.findOne({name: newCart.name, handleLength:req.body.handleLength, bucketPlacement:req.body.bucketPlacement,
      bucketSize:req.body.bucketSize, wheelSize:req.body.wheelSize},
      addTripCallback(response, newCart, totalCarts));
  };
}

function addCartData(req, response) {
  var newCart = new UserCart(req.body);
  newCart.name = req.params.name;

  // get the count of the number of carts saved for the username (req.params.name)
  //     and pass on a newly constructed UserCart populated by the user's parameters
  //     to the call back method so that the method can provide the actual save method
  //     with both the number of carts this user has saved and this newly constructred
  //     cart if there isn't already a previously saved cart that matches the user's
  //     parameters if there was one.
  UserCart.count({name: req.params.name}, addTripAfterCountCallBack(req, response, newCart));
}

//For Testing only
//curl -i -X POST -H 'Content-Type: application/json' http://localhost:3000/usercart/erik/removeCartData
function removeCartData(req, response) {
  var name = req.params.name;
  UserCart.remove({name: name}, function (err) {
    if (err) {
      response.send(500);
      return;
    }
    response.send(404);
    return;
  });
}


function findAllCarts(req, response) {
  UserCart.find({name: req.params.name}, function(error, carts){
    if (error) {
      response.send(500);
      return;
    }
    response.send(200, filterHidden(carts));
    return;
  });
}

exports.usercart = {
  findAllCarts: findAllCarts,
  addCartData: addCartData,
  removeCartData: removeCartData,
};


// userSchema represents a user account
var userSchema = mongoose.Schema({
  name: {type: String, required: true},
});

var User = mongoose.model('User', userSchema);


function createUser(req, response) {
  var name = req.body.name;
  User.findOne({name: name}, function(error, user) {
    if (error) {
      response.send(500);
      return;
    }

    if (user) {
      response.send(500);
      return;
    }

    var user = new User(req.body);
    user.save(saveCallback(response));
  });
}

function findUserByName(req, response) {
  User.findOne({name: req.params.name}, function(error, user) {
    if (error) {
      response.send(500);
      return;
    }

    if (!user) {
      response.send(404);
      return;
    }

    response.send(200, filterHidden(user));
  });
}

exports.users = {
  create: createUser,
  findByName: findUserByName,
};


//Helper Functions
function saveCallback(response) {
  return function(err, result) {
    if(err) {
     console.log(err);
     response.send(500);
    } else {
     console.log(result);
     response.send(200, filterHidden(result));
    }
  };
}

function filterHidden(v) {
  return JSON.stringify(v, function(n, v) {
    if (n == '_id' || n == '__v')
      return undefined;
    return v;
  })
}


















var userChallengeSchema = mongoose.Schema({
  name: {type: String, required: true}, //username
  openEndedQuestions: [{ question: String, answer: String, time: Date }],
  feature: {type: String, enum: ['handleLength', 'wheelSize', 'bucketSize', 'bucketPlacement']},
  state: {type: String, enum: ['selectAFeature', 'experiment', 'claim', 'emptyNotebook',
                               'knowForSure', 'selectRecords', 'singleCase', 'multiCases'],
          default:'selectAFeature'},
  challengeID: {type: Number, required: true},
});

var UserChallenge = mongoose.model('UserChallenge', userChallengeSchema);

function updateStateCallback(response, newChallenge, totalChallenges) {
  return function(error, challenge) {
    if (challenge) {
      challenge.state = newChallenge.state;
      challenge.save(saveCallback(response));
    } else {
      newChallenge.challengeID = totalChallenges + 1; //TODO not really save to use count as ID because of the DB calls are async
      newChallenge.save(saveCallback(response));
    }
  };
}

function addChallengeAfterCountCallBack(req, response, newChallenge) {
  return function(error, totalChallenges) {
    console.log(totalChallenges);

    UserChallenge.findOne({name: newChallenge.name, challengeID:req.body.challengeID},
      updateStateCallback(response, newChallenge, totalChallenges));
  };
}

function addChallenge(req, response) {
  var newChallenge = new UserChallenge(req.body);
  newChallenge.name = req.params.name;

  // get the count of the number of challenges saved for the username (req.params.name)
  //     and pass on a newly constructed UserChallenge populated by the user's parameters
  //     to the call back method so that the method can provide the actual save method
  //     with both the number of challenges this user had and this newly constructred
  //     challenge if there isn't already a previously saved challenge with this id.
  UserChallenge.count({name: req.params.name}, addChallengeAfterCountCallBack(req, response, newChallenge));
}

//For Testing only
//curl -i -X POST -H 'Content-Type: application/json' http://localhost:3000/usercart/erik/removeChallenges
function removeChallenges(req, response) {
  var name = req.params.name;
  UserChallenge.remove({name: name}, function (err) {
    if (err) {
      response.send(500);
      return;
    }
    response.send(404);
    return;
  });
}

function findAllChallenges(req, response) {
  UserChallenge.find({name: req.params.name}, function(error, challenges){
    if (error) {
      response.send(500);
      return;
    }
    response.send(200, filterHidden(challenges));
    return;
  });
}

function findChallengeByID(req, response) {
  UserChallenge.findOne({name: req.params.name, challengeID: req.params.challengeID}, function(error, challenge){
    if (error) {
      response.send(500);
      return;
    }

    response.send(200, filterHidden(challenge));
    return;
  });
}


exports.userchallenge = {
  findAllChallenges: findAllChallenges,
  findChallengeByID: findChallengeByID,
  addChallenge: addChallenge,
  removeChallenges: removeChallenges,
};

