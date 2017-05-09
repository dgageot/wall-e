"use strict";

var lab = angular.module('wall', []);

lab.controller('BuildWallCtrl', ($scope, $http, $interval, $q) => {
    function refresh() {
      var deferred = $q.defer();
      var urlCalls = [];

      $http.get('/github/repos/docker/pinata/pulls').success(data => {
        data.forEach(pr => urlCalls.push($http.get(statusesUrl(pr))));

        $q.all(urlCalls).then(results => {
          for (var i = 0; i < results.length; i++) {
            data[i].status = results[i].data;
          }
          $scope.pulls = data;
        });
      });

      $http.get('/jenkins/api/json?tree=jobs[name,builds[building,number,result,runs[building,number,result]]{0,5}]').success(function(data) {
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

var statusesUrl = pr => {
  return pr.statuses_url.replace('https://api.github.com', '/github').replace('/statuses', '/status');
}
