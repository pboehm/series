package streams

var ServerStaticHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Series - Streams</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">

    <style>
        .navbar-fixed {
            height: 101px;
        }

        .navbar-fixed .actions {
            display: flex;
            justify-content: space-between;
            padding: 10px 20px;
        }

        .navbar-fixed .actions .hidden {
            display: none;
        }

        nav.nav-extended .nav-wrapper {
            min-height: 45px !important;
        }

        nav .brand-logo {
            line-height: 40px !important;
        }

        .container {
            padding-top: 10px;
            width: 95%;
        }

        .container .collapsible .collapsible-header {
            align-items: center;
        }

        .container .collapsible .collapsible-header .series-title {
            display: flex;
            flex-grow: 1;
            font-weight: bold;
        }

        .container .collapsible .collapsible-header .episode-count {
            padding-left: 10px;
        }

        .container .collapsible .collapsible-body {
            padding: 0;
        }

        .container .collection .collection-item {
            padding: 15px;
        }

        .container .collection-item .top {
            display: flex;
            justify-content: space-between;
            align-items: normal;
        }

        .container .collection-item .top .action span {
            padding-left: 20px !important;
            height: 15px !important;
        }

        .container .collection-item .bottom {
            padding-top: 10px;
        }

        .container .collection-item .bottom .buttons button {
            margin-top: 10px;
        }

    </style>
</head>

<body>
<div class="navbar-fixed">
    <nav class="nav-extended teal lighten-1" role="navigation">
        <div class="nav-wrapper container"><a id="logo-container" href="#" class="brand-logo">Streams</a></div>
        <div class="nav-content">
            <div class="actions">
                <a class="orange waves-effect waves-light btn hidden" id="load-button">Load</a>
                <a class="grey waves-effect waves-light btn" id="refresh-button">Refresh</a>
                <a class="orange waves-effect waves-light btn disabled" id="mark-watched-button">Mark as watched</a>
            </div>
        </div>
    </nav>
</div>

<div class="container" id="series-container">

</div>

<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.1/mustache.min.js"></script>

<script id="template-series" type="x-tmpl-mustache">
    [[#groups]]
    <ul class="collapsible">
        <li>
            <div class="collapsible-header">
                <div class="series-title"><i class="material-icons">live_tv</i>[[ series ]]</div>
                <div class="episode-count">[[ episodes.length ]]</div>
            </div>
            <div class="collapsible-body">

                <ul class="collection">
                    [[#episodes]]
                    <li class="collection-item">
                        <div class="top">
                            <div class="title">[[ filename ]]</div>
                            <div class="action">
                                 <label>
                                  <input type="checkbox" class="watched-checkbox filled-in" [[#watched]]checked[[/watched]] name="[[ id ]]"/>
                                  <span></span>
                                </label>
                            </div>
                        </div>
                        <div class="bottom">
                            <div class="hoster">
                                <select id="select-[[ __id ]]">
                                    [[#links]]
                                    <option value="[[ id ]]" data-link="[[ link ]]">[[ hoster ]]</option>
                                    [[/links]]
                                </select>
                            </div>
                            <div class="buttons">
                                <button data-episode-id="[[ __id ]]"
                                    class="link-button grey waves-effect waves-light btn-small">Open</button>

                                [[#linkActions]]
                                    <button class="link-action-button grey waves-effect waves-light btn-small"
                                        data-episode-id="[[ __id ]]" data-action="[[ id ]]">[[ title ]]</button>
                                [[/linkActions]]
                            </div>
                        </div>
                    </li>
                    [[/episodes]]
                </ul>
            </div>
        </li>
    </ul>
    [[/groups]]
</script>

<script>
    var watchedEpisodes = (!!window.localStorage) ? window.localStorage : {};
    var appState = {
        groups: [],
        linkActions: [],
        globalActions: []
    };

    function hasEpisodeBeenWatched(episodeId) {
        return !!watchedEpisodes[episodeId];
    }

    function markEpisodeAsWatched(episodeId) {
        watchedEpisodes[episodeId] = "1";
    }

    function unmarkEpisodeAsWatched(episodeId) {
        delete watchedEpisodes[episodeId];
    }

    function manageMarkWatchedButton() {
        var checkedCount = $(".watched-checkbox:checked").length;

        if (checkedCount > 0) {
            $("#mark-watched-button").removeClass("disabled");
        } else {
            $("#mark-watched-button").addClass("disabled");
        }
    }

    function callLinkAction(button, action, linkId) {
        var originalText = button.textContent;

        button.text = originalText.replace(/./g, ".");
        button.classList.add("disabled");

        fetch("/api/actions/link/" + action + "/" + linkId, {"method": "POST"})
            .then(function (response) {
                return response.json();
            })
            .then(function (success) {
                console.log(success);
                button.textContent = originalText;
                button.classList.remove("disabled");
            }, function (error) {
                console.log(error);
                button.textContent = originalText;
                button.classList.remove("disabled");
            });
    }

    function registerHandlers() {
        $(".watched-checkbox").click(function(e) {
            if (this.checked) {
                markEpisodeAsWatched(this.name);
            } else {
                unmarkEpisodeAsWatched(this.name);
            }

            manageMarkWatchedButton();
        });

        $(".link-button").click(function(e) {
            var episodeId = this.dataset.episodeId;
            markEpisodeAsWatched(episodeId);
            $(".watched-checkbox[name=" + episodeId + "]").attr('checked', true);
            manageMarkWatchedButton();

            var link = $("#select-" + episodeId + " :selected").data("link");
            window.open(link, '_blank');
        });

        $(".link-action-button").click(function(e) {
            var episodeId = this.dataset.episodeId;
            markEpisodeAsWatched(episodeId);
            $(".watched-checkbox[name=" + episodeId + "]").attr('checked', true);
            manageMarkWatchedButton();

            var action = this.dataset.action;
            var linkId = $("#select-" + episodeId + " :selected").attr("value");
            callLinkAction(this, action, linkId);
        });
    }

    function renderState() {
        var template = $('#template-series').html();
        Mustache.parse(template);
        var rendered = Mustache.render(template, appState, null, ['[[', ']]']);
        $('#series-container').html(rendered);

        registerHandlers();
        manageMarkWatchedButton();

        $('.collapsible').collapsible();
        $('select').formSelect();
    }

    function loadActions() {
        var linked = fetch("/api/actions/link").then(function (response) { return response.json() });
        var global = fetch("/api/actions/global").then(function (response) { return response.json() });

        Promise.all([linked, global]).then(function (successes) {
            appState.linkActions = successes[0];
            appState.globalActions = successes[1];
            renderState();
        }, function (error) {
            console.log(error);
        });
    }

    function loadLinks() {
        fetch("/api/links/grouped").then(function (response) {
            return response.json();
        }).then(function (success) {
            var refreshButton = $("#refresh-button");
            var loadButton = $("#load-button");
            if (!!success.ready) {
                refreshButton.removeClass("hidden");
                loadButton.addClass("hidden");
            } else {
                loadButton.removeClass("hidden");
                refreshButton.addClass("hidden");
            }

            var groups = success.links;
            groups.forEach(function (group) {
                var episodes = group["episodes"];
                episodes.forEach(function (episode) {
                    episode["__id"] = episode["id"];
                    episode["watched"] = hasEpisodeBeenWatched(episode["id"]);

                    episode["links"].forEach(function (link) {
                        link["episodeId"] = episode["id"];
                    });
                });
            });

            appState.groups = groups;

            renderState();
        }, function (error) {
            console.log(error);
        });
    }

    function refreshInBackend(button) {
        var originalText = button.text;

        button.text = "...";
        button.classList.add("disabled");

        fetch("/api/links/refresh", {"method": "POST"}).then(function (response) {
            return response.json();
        }).then(function (success) {
            button.text = originalText;
            button.classList.remove("disabled");
            loadLinks();
        }, function (error) {
            console.log(error);
            button.text = originalText;
            button.classList.remove("disabled");
        });
    }

    function markAsWatchedInBackend(button) {
        var originalText = button.text;

        button.text = "...";
        button.classList.add("disabled");

        var episodeIds = [];

        $(".watched-checkbox:checked").each(function () {
            episodeIds.push(this.name);
        });

        fetch("/api/links/watched", {"method": "POST", "body": JSON.stringify(episodeIds)}).then(function (response) {
            return response.json();
        }).then(function (success) {
            button.text = originalText;
            button.classList.remove("disabled");
            loadLinks();
        }, function (error) {
            console.log(error);
            button.text = originalText;
            button.classList.remove("disabled");
        });
    }

    (function () {
        $("#load-button").click(function() {
            loadLinks();
        });

        $("#refresh-button").click(function() {
            refreshInBackend(this);
        });

        $("#mark-watched-button").click(function() {
            markAsWatchedInBackend(this);
        });

        loadActions();
        loadLinks();
    })();
</script>
</body>
</html>
`
