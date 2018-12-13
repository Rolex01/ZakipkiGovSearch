var app = angular.module('adminPanel', []);

app.controller('UsersListCtrl', function ($scope, $http) {
  $scope.loadData = function() {
	  $http.post('/control/a/users/list').success(function(data) {
		$scope.users = data;
		$scope.users.map(function(u){
			u.createdDate = new Date(u.created)
			m = moment(u.createdDate)
			u.created=m.format("DD.MM.YYYY")
			u.endDate = new Date(u.lifetime)
			m = moment(u.endDate)
			u.lifetime=m.format("DD.MM.YYYY")
			switch(u.permission) {				
				case -2:
					u.accessGroup = "Клиент (Просрочена оплата)"
					break
				case -1:
					u.accessGroup = "Демо (Просрочен)"
					break
				case 0:
					u.accessGroup = "Демо"
					break
				case 1:
					u.accessGroup = "Клиент"
					break
				case 2:
					u.accessGroup = "Сотрудник"
					break
				case 90:
					u.accessGroup = "Администратор"
					break
				case 100:
					u.accessGroup = "Разработчик"
					break
				default:
					u.accessGroup = "Не указан"
					break
				}		   
		})
	    console.log($scope.users)	
	  })	  
	  $scope.orderProp = 'name';  
	  $scope.orderByNameFlag = false
	  $scope.orderByEmailFlag = false
	  $scope.orderByPermissionFlag = false
	  $scope.orderByCreatedFlag = false
	  $scope.orderByLifetimeFlag = false
	  $scope.orderBySLimitFlag = false
	  $scope.orderByRLimitFlag = false
  }
  $scope.orderByName = function() {
	if ($scope.orderByNameFlag == true) {
		$scope.orderProp = 'name';
		$scope.orderByNameFlag = false

	} else {
		$scope.orderProp = '-name';
		$scope.orderByNameFlag = true

	}
  }

  $scope.orderByEmail = function() {
	if ($scope.orderByEmailFlag == true) {
		$scope.orderProp = 'email';
		$scope.orderByEmailFlag = false

	} else {
		$scope.orderProp = '-email';
		$scope.orderByEmailFlag = true

	}
  }

  $scope.orderByPermission = function() {
	if ($scope.orderByPermissionFlag == true) {
		$scope.orderProp = 'permission';
		$scope.orderByPermissionFlag = false

	} else {
		$scope.orderProp = '-permission';
		$scope.orderByPermissionFlag = true
	}
  }

  $scope.orderByCreated = function() {
	if ($scope.orderByCreatedFlag == true) {
		$scope.orderProp = 'startDate';
		$scope.orderByCreatedFlag = false
	} else {
		$scope.orderProp = '-startDate';
		$scope.orderByCreatedFlag = true
	}
  }

  $scope.orderByLifetime = function() {
	if ($scope.orderByLifetimeFlag == true) {
		$scope.orderProp = 'endDate';
		$scope.orderByLifetimeFlag = false
	} else {
		$scope.orderProp = '-endDate';
		$scope.orderByLifetimeFlag = true
	}
  }

  $scope.orderBySLimit = function() {
	if ($scope.orderBySLimitFlag == true) {
		$scope.orderProp = 'sessionslimit';
		$scope.orderBySLimitFlag = false
	} else {
		$scope.orderProp = '-sessionslimit';
		$scope.orderBySLimitFlag = true
	}
  }

  $scope.orderByRLimit = function() {
	if ($scope.orderByRLimitFlag == true) {
		$scope.orderProp = 'requestRLimit';
		$scope.orderByRLimitFlag = false
	} else {
		$scope.orderProp = '-requestRLimit';
		$scope.orderByRLimitFlag = true
	}
  }

  $scope.loadData()
  

});