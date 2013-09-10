

function WebhookListCtrl($scope, webhookId, $cookies, $timeout) {
    $scope.recentWebhooks = {};

    if(typeof $cookies[webhookId] == "undefined") {
        $cookies['wh-' + webhookId] = new Date().toLocaleString();
    }

    $scope.$watch(function() {
        var hash = 0;
        for(var webhookId in $cookies) {
            hash++;
        }

        return hash;
    }, function() {
        $scope.recentWebhooks = [];
        for (var webhookId in $cookies) {
            var webhookIdMatch = webhookId.match(/wh-(\w+)/);
            if(webhookIdMatch) {
                var webHook = {
                    id: webhookIdMatch[1],
                    date: $cookies[webhookIdMatch[0]]
                };
                $scope.recentWebhooks.push(webHook);
            }
        }
        $scope.recentWebhooks.sort(function(a, b){
            return Date.parse(b.date) - Date.parse(a.date);
        });
    });

    $scope.clearRecent = function(){
        for (var webhookId in $cookies) {
            var webhookIdMatch = webhookId.match(/wh-(\w+)/);
            if(webhookIdMatch) {
                delete $cookies[webhookIdMatch[0]];
            }
        }
    }
}

function RequestListCtrl($scope, webhookId, staticUrl, socket) {
    $scope.requests = [];
    $scope.totalSubscribers = 0;
    $scope.enableSounds = true;

    $scope.supportsAudioCodec = function(codec) {
        var a = document.createElement('audio');
        return !!(a.canPlayType && a.canPlayType('audio/mpeg;').replace(/no/, ''));
    }

    $scope.audioTagForUrl = function(audioUrl, preload){
        var a = document.createElement('audio');
        a.preload = typeof preload == 'undefined' ? false: preload;

        if($scope.supportsAudioCodec('audio/ogg')) {
            a.src = audioUrl + '.ogg';
        } else if($scope.supportsAudioCodec('audio/mpeg')) {
            a.src = audioUrl;
        }

        return a;
    };

    var effectsPath = staticUrl + 'audio/effect/';
    $scope.soundEffects = {
        normal: $scope.audioTagForUrl(effectsPath + 'normal.mp3', true),
        quiet: $scope.audioTagForUrl(effectsPath + 'quiet.mp3'),
        medium: $scope.audioTagForUrl(effectsPath + 'medium.mp3'),
        loud: $scope.audioTagForUrl(effectsPath + 'loud.mp3'),
        error: $scope.audioTagForUrl(effectsPath + 'error.mp3')
    };

    socket.on('connect', function(){
        socket.emit('join', webhookId);
    });

    socket.on('reconnect', function () {
        console.info('Reconnected to the server');
    });

    socket.on('reconnecting', function () {
        console.info('Attempting to re-connect to the server');
    });

    socket.on('error', function (e) {
        console.error(e ? e : 'A unknown error occurred');
    });

    socket.on('new_request', function(message) {
        message.date = new Date();
        message.dateNumber = Number(message.date);
        message.hasInfo = false;
        message.hasHeaders = false;
        message.hasPost = message.post != 'null' && message.post != '';
        message.hasGet = message.get != 'null' && message.get != '';
        message.hasFiles = message.files != 'null' && message.files != '';
        $scope.requests.push(message);

        $scope.playSoundEffect(message.soundEffect);
    });

    socket.on('subscriber_joined', function(message){
        $scope.totalSubscribers = message.totalSubscribers;
    });

    $scope.playSoundEffect = function(soundEffect) {
        if($scope.enableSounds) {
            $scope.soundEffects[soundEffect].play();
        }
    }

    $scope.removeRequest = function(index){
        $scope.requests.splice(index, 1);
    }
}


app.directive('snippet', ['$timeout', '$interpolate', function($timeout, $interpolate) {
    function replaceURLWithHTMLLinks(text) {
        var exp = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
        return text.replace(exp, '<a href="$1" target="_blank">$1</a>');
    }

    return {
        restrict: 'E',
        template:'<pre><code class="prettyprint" ng-transclude></code></pre>',
        replace:true,
        transclude:true,
        link:function(scope, elm, attrs){             
            var tmp =  $interpolate(elm.find('code').text())(scope);
             $timeout(function() {                
                elm.find('code').html(replaceURLWithHTMLLinks(hljs.highlightAuto(tmp).value));
              }, 0);
        }
    };
}]);

app.filter('reverse', function() {
    return function(items) {
        return items.slice().reverse();
    };
});



app.directive('timeSince', function($timeout, dateFilter) {
    // return the directive link function. (compile function not needed)
    return function(scope, element, attrs) {
        var timeoutId; // timeoutId, so that we can cancel the time updates
        var since;   

        function timeSinceFriendly(date) {
            var seconds = Math.floor((new Date() - date) / 1000);

            var interval = Math.floor(seconds / 31536000);

            if (interval > 1) {
                return interval + " years";
            }
            interval = Math.floor(seconds / 2592000);
            if (interval > 1) {
                return interval + " months";
            }
            interval = Math.floor(seconds / 86400);
            if (interval > 1) {
                return interval + " days";
            }
            interval = Math.floor(seconds / 3600);
            if (interval > 1) {
                return interval + " hours";
            }
            interval = Math.floor(seconds / 60);
            if (interval > 1) {
                return interval + " minutes";
            }
            return Math.floor(seconds) + " seconds";
        }

        // used to update the UI
        function updateTime() {
            element.text(timeSinceFriendly(since));  
        }

        scope.$watch(attrs.timeSince, function(value) {
            since = new Date(value);
            updateTime();
        });

        // schedule update in one second
        function updateLater() {
            // save the timeoutId for canceling
            var timeoutId = $timeout(function() {
                updateTime(); // update DOM
                updateLater(); // schedule another update
            }, 1000);
        }

        // listen on DOM destroy (removal) event, and cancel the next UI update
        // to prevent updating time ofter the DOM element was removed.
        element.bind('$destroy', function() {
            $timeout.cancel(timeoutId);
        });

        updateLater(); // kick off the UI update process.
    }
});
