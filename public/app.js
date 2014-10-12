'use strict';

// Declare app level module which depends on views, and components
angular.module('myApp', [
  'ngRoute',
])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider
  .when('/', {
    templateUrl: 'static/views/main.html',
    controller: 'MainCtrl'
  })
  .when('/game/:gameid', {
    templateUrl: 'static/views/game.html',
    controller: 'GameCtrl'
  })
  .otherwise({redirectTo: '/'});
}])

.controller('MainCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	$scope.createGame = function(gameType) {
		if (gameType === 1) {
			$http.post('/game', {}).then(function success(response) {
				var runnerId = response.data.id;
				$location.path('/game/' + runnerId);
			})
		} else {
			$http.post('/game/ivan', {}).then(function success(response) {
				var runnerId = response.data.id;
				$location.path('/game/' + runnerId);
			})
		}
		
		
	}
}])

.controller('GameCtrl', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location) {
	$scope.message = "Waiting for opponent (Share the url with a friend)";
	$scope.player = {};
	$scope.opponent = {};
	$scope.board = ["","","","","","","","",""];;
	$scope.myTurn = false;
	$scope.gameOver = false;

	var source;

	function handleNewPlayerEvent(event) {
		var player = event.data.split(":");
		$scope.$apply(function() {
			if (!$scope.player.id) {
				$scope.player.id = player[0];
				$scope.player.mark = player[1];
				if (player[0] == "1") {
					$scope.opponent.id = "0";
					$scope.opponent.mark = "X";
				}
			} else {
				$scope.opponent.id = player[0];
				$scope.opponent.mark = player[1];
			}
		});
	}

	function handleGameReadyEvent(event) {
		$scope.gameOver = false;
		var playerId = event.data
		$scope.$apply(function() {
			$scope.message = "";
			$scope.board = ["","","","","","","","",""];
			$scope.myTurn = $scope.player.id == playerId;
		});
	} 

	function handleMoveEvent(event) {
		var data = event.data.split(":"),
			playerId = data[0],
			place = parseInt(data[1]);
		if (playerId != $scope.player.id) {
			$scope.$apply(function() {
				$scope.board[place] = $scope.opponent.mark;
				$scope.myTurn = true;
			});
		} else {
			// move is already made
		}
	}

	function handleGameOverEvent(event) {
		var winnerMark = event.data;
		$scope.$apply(function() {
			$scope.gameOver = true;
			if ($scope.player.mark === winnerMark) {
				$scope.draw = false;
				$scope.won = true;
			} else if ($scope.opponent.mark === winnerMark){
				$scope.draw = false;
				$scope.won = false;
			} else {
				$scope.draw = true;
			}
		});
	}

	function handleRestartEvent(event) {
		console.log(event);
	}

	$scope.restart = function() {
		$http({
			method: 'PATCH',
			url: '/game/'+$routeParams.gameid+'/restart', 
			data: angular.toJson({
			playerId: parseInt($scope.player.id),
		})}).then(function success(response) {
			$scope.message = "Waiting for opponent to accept"
		})
	}

	function init() {
		$scope.loading = true;
		// make sure the game exists
		$http.get('/game/'+$routeParams.gameid).then(
			function success(response) {
				$scope.loading = false;
				//register to event source
				source = new EventSource('/game/'+$routeParams.gameid+'/updates');
				// register callbacks
				source.addEventListener("newPlayer", handleNewPlayerEvent);
				source.addEventListener("gameReady", handleGameReadyEvent);
				source.addEventListener("move", handleMoveEvent);
				source.addEventListener("gameOver", handleGameOverEvent);
				source.addEventListener("restart", handleRestartEvent);
			},
			function error(response) {
				$scope.loading = false;
				$scope.message = "Game does not exist";
				window.setTimeout(function() {
					$location.path("/")
				}, 1500);
			});
	}
	
	init();


	$scope.play = function(place) {
		if ($scope.myTurn && $scope.board[place] === ""){
			$http({
				method: 'PATCH',
				url: '/game/'+$routeParams.gameid+'/move', 
				data: angular.toJson({
				playerId: parseInt($scope.player.id),
				place: place
			})}).then(function success(response) {
				if (response.data.status == "success") {
					$scope.board[place] = $scope.player.mark;
					$scope.myTurn = false;
				}
			})
		}
	}

	$scope.getTileClass = function(place) {
		if ($scope.board[place] === "X") {
			return "fa-times";
		} else if ($scope.board[place] === "O") {
			return "fa-circle-o";
		}
		return;
	}
}]);