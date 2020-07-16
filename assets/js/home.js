var loading = false;
// initialize handlebars templates
var templates = {};
$("script[id^=template]").each(function(){
  var name = $(this).attr("id").substring(9);
  var source   = $(this).html();
  templates[name] = Handlebars.compile(source);
});

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
        window.location.href = data.redirect;
    }
  }).always(function(){
      loading = false;
      $("a.refresh-button").removeClass("is-loading");
      $("#check_form .button").removeClass("is-loading");
  });
  return false;
};

// on ready
$(function(){

  // handle form submission
  $("form#check_form").submit(loadData);
});