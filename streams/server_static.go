package streams

var ServerStaticHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
  <title>Series - Streams</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
</head>

<body>
  <nav class="teal lighten-1" role="navigation">
    <div class="nav-wrapper container"><a id="logo-container" href="#" class="brand-logo">Streams</a></div>
  </nav>

  <div class="container" id="content">

  </div>

  <script src="http://code.jquery.com/jquery-3.3.1.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.1/mustache.min.js"></script>

  <script id="template-series" type="x-tmpl-mustache">
    [[#groups]]
    <ul class="collection with-header">
        <li class="collection-header"><h4>[[ series ]]</h4></li>
        [[#episodes]]
        <li class="collection-item">
            <div>
                [[ filename ]]
                <a href="#" class="secondary-content"><i class="material-icons">check_box</i></a>
            </div>
            <div style="overflow-x: scroll; white-space: nowrap; padding-top: 10px;">
                [[#links]]
                <a href="[[ link ]]" target="_blank" class="waves-effect waves-light btn-small">[[ hoster ]]</a>
                [[/links]]
            </div>
        </li>
        [[/episodes]]
    </ul>
    [[/groups]]
  </script>

  <script>
    function renderLinks(groups) {
      var template = $('#template-series').html();
      Mustache.parse(template);   // optional, speeds up future uses
      var rendered = Mustache.render(template, {groups: groups}, null, ['[[', ']]']);
      $('#content').html(rendered);
    }
    
    (function() {
        fetch("/api/links/grouped").then(function (response) { 
            return response.json();
         }).then(function (success) { 
             var links = success.links;
             
             var groups = [];
             
             Object.keys(links).forEach(function (groupName) { 
                 groups.push({
                    series: groupName,
                    episodes: links[groupName]
                 });
             });
             
             console.log(groups);
             renderLinks(groups)
         }, function (error) { 
             console.log(error);
         });
    })();
  </script>
</body>
</html>
`
