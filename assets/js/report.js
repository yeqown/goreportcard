Handlebars.registerHelper('gradeMessage', function(grade, options) {
  var gradeMessages = {
    "A+": "Excellent!",
    "A": "Great!",
    "B": "Not bad!",
    "C": "Needs some work",
    "D": "Needs lots of improvement",
    "E": "Urgent improvement needed",
    "F": "... is for lots of things to Fix!"
  };
  return gradeMessages[grade];
});

// add a helper for picking the progress bar colors
Handlebars.registerHelper('color', function(percentage, options) {
  switch(true){
    case percentage < 30:
      return 'is-danger';
    case percentage < 50:
      return 'is-warning';
    case percentage < 80:
      return 'is-info';
    default:
      return 'is-success';
  };
});

Handlebars.registerHelper('isfalse', function(percentage, options) {
  return percentage == false;
});

var allowedLinkDomains = ["https://github.com/", "https://bitbucket.org/",
  "https://golang.org/", "https://go.googlesource.com/"];

// initialize handlebars templates
var templates = {};
$("script[id^=template]").each(function(){
  var name = $(this).attr("id").substring(9);
  var source   = $(this).html();
  templates[name] = Handlebars.compile(source);
});

var shrinkHeader = function(){
  var $hero = $("section.hero");
  $hero.slideUp();
}

var populateResults = function(data, domain){
    var checks = data.scores;
    var $resultsText = $(".results-text");
    var $resultsDetails = $('.results-details').empty();

    for (var i = 0; i < allowedLinkDomains.length; i++) {
      if (data.resolvedRepo.indexOf(allowedLinkDomains[i]) == 0) {
        data.link = data.resolvedRepo;
      }
    }
    data.use_an = data.grade == "A" || data.grade == "A+";
    data.grade_encoded = encodeURIComponent(data.grade);
    $resultsText.html($(templates.grade(data)));
    var $table = $(".results");
    $table.html('<p class="panel-heading">Results</p>');
    for (var i = 0; i < checks.length; i++) {
        checks[i].percentage = parseInt(checks[i].percentage * 100.0);
        var $headRow = $(templates.check(checks[i]));
        $headRow.on("click", function(){
        $(this).closest("nav").find(".is-active").removeClass("is-active");
          $(this).toggleClass("is-active");
        });
        $headRow.appendTo($table);
        if (i == 0) {
            $headRow.toggleClass("is-active");
        }

        var $details = $(templates.details(checks[i]));
        $details.appendTo($resultsDetails);
    }
    $(".container-suggestions").addClass('hidden');
    $(".container-results").removeClass('hidden').slideDown();

    $lastRefresh = $(templates.lastrefresh(data));
    $div = $(".container-update").html($lastRefresh);
    $div.find("a.refresh-button").on("click", function(e){
      loadData.call($("form#check_form")[0], false);
      $(this).addClass('is-loading');
      return false;
    });

    var badgeData = {
        url: "https://" + domain + "/report/" + data.repo,
        image_url: "https://" + domain + "/badge/" + data.repo,
    }
    var $badgeDropdown = $(templates.badgedropdown(badgeData));
    $badgeDropdown.find("input").on("click", function(){
        $(this).select();
    });
    $(".badge-col").append($badgeDropdown);
    $(".badge-col .badge").on("click", function(){
        $(this).closest(".badge-col").find("#badge_dropdown").toggleClass("hidden");
    });
};

function alertMessage(msg){
  var html = templates.alert({message: msg});
  var $alert = $(html);
  $alert.find(".delete").on("click", function(){
      $(this).closest(".notification").remove();
  });
  $("#notifications").children().remove();
  $alert.hide();
  $alert.appendTo("#notifications");
  $alert.slideDown();
}

var loadData = function(getRequest){
  loading = true;
  var $form = $(this),
      url = $form.attr("action"),
      method = $form.attr("method"),
      data = {};
    $form.serializeArray().map(function(x){data[x.name] = x.value;});

    if(!data["repo"]) {
        alertMessage("Input cannot be empty. Please enter a valid repository path");
        return false;
    }

    $("#check_form .button").addClass("is-loading");
  $.ajax({
      type: getRequest ? "GET" : "POST",
      url: url,
      data: data,
      dataType: "json"
  }).fail(function(xhr, status, err){
      alertMessage("There was an error processing your request: " + xhr.responseText);
  }).done(function(data, textStatus, jqXHR){
      if (data.redirect) {
          location.replace(data.redirect);
      }
  }).always(function(){
      loading = false;
      $("a.refresh-button").removeClass("is-loading");
      $("#check_form .button").removeClass("is-loading");
      $(".container-loading").slideUp();
  });
  return false;
};

var hideResults = function(){
  $(".container-results").hide();
};

// on ready
$(function(){

  if (loading) {
      // we need to load the results
      loadData.call($("form#check_form")[0], true);
  } else {
      populateResults(JSON.parse(response), domain);
      $(".container-loading").slideUp();
  }

  // handle form submission
  $("form#check_form").submit(loadData);

  // sticky menu
  $(window).scroll(function() {
      if ($(this).scrollTop() >= 240) {
          $('nav.results').addClass('stickytop');
      }
      else {
          $('nav.results').removeClass('stickytop');
      }
  });
});
