"use strict";

var lab = angular.module('wall', []);

lab.controller('BuildWallCtrl', ($scope, $http, $interval, $q) => {
    function refresh() {
      var deferred = $q.defer();
      var urlCalls = [];

      $http.get('/proxy/pulls').success(data => {
        data.forEach(pr => urlCalls.push($http.get("/proxy/status/" + pr.head.sha)));

        $q.all(urlCalls).then(results => {
          for (var i = 0; i < results.length; i++) {
            data[i].status = results[i].data;
            data[i].status.statuses = data[i].status.statuses.sort((l, r) => l.context.localeCompare(r.context));
          }
          $scope.pulls = data;
        });
      });

      $http.get('/proxy/jobs').success(function(data) {
        $scope.macMaster = jobPlus(data.jobs.find(job => job.name == 'pinata-macos-master'));
        $scope.winMaster = jobPlus(data.jobs.find(job => job.name == 'pinata-win-master'));
      });
    }

    refresh();
    $interval(_ => refresh(), 5000)
});

var lastResult = job => {
  var lastJobDone = job.builds.find(build => build.result);
  return (lastJobDone && lastJobDone.result) || "UNKNOWN";
};

var jobPlus = job => {
  job.result = lastResult(job);
  job.builds = job.builds.map(build => {
    build.runs = build.runs.filter(run => run.number == build.number);
    build.ok = build.runs.filter(run => run.result == 'SUCCESS');
    build.ko = build.runs.filter(run => run.result == 'FAILURE' || run.result == 'ABORTED');
    return build;
  });
  return job;
}
