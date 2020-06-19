const search_input = document.getElementById("search-input");
const search_image = document.getElementById("search-image");
const more_image = document.getElementById("more-image");

var xhttp = new XMLHttpRequest();
var card_num = 24;

search_input.addEventListener('keyup', InputHandler);
search_image.addEventListener('click', ImageHandler);
more_image.addEventListener('click', MoreHandler);

function getBaseUrl() {
    return window.location.href.match(/^.*\//);
}

xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
        AddFilms(this);
    }
};

function AddFilms (xml) {
    var x = xml.responseXML.getElementsByTagName("card");
    var release_date;
    var i;
    for (i = 0; i < x.length; i++){
        if (x[i].getElementsByTagName("text")[0].childNodes.length != 0) {
            release_date = x[i].getElementsByTagName("text")[0].childNodes[0].nodeValue;
        } else {
            release_date = "";
        }
        AddFilm(
            x[i].getAttribute("watchable"),
            x[i].getElementsByTagName("id")[0].childNodes[0].nodeValue,
            x[i].getElementsByTagName("picture")[0].childNodes[0].nodeValue,
            x[i].getElementsByTagName("title")[0].childNodes[0].nodeValue,
            release_date,
        );
    }
}

function AddFilm (watchable, id, poster, title, releaseDate) {
    var new_film;

    if (watchable == "true") {
        new_film = '<a class="watchable_film_item" HREF="watch?v=' + id + '">';
    }
    else {
        new_film = '<a class="non-watchable_film_item" HREF="watch?v=' + id + '">';
    }

    if (poster != "") {
        new_film += '<img src="' + poster + '">';
    }
    new_film += '<div>' +
    '<p><b>' + title + '</b></p>' +
    '<p>' + releaseDate + '</p>' +
    '</div></a>';

    document.getElementById("results").innerHTML += new_film;
}

function InputHandler (key_event) {
    if (event.keyCode === 13) {
        window.location.replace("?q=" + encodeURI(search_input.value));
    }
}

function ImageHandler () {
    window.location.replace("?q=" + encodeURI(search_input.value));
}

function MoreHandler () {
    var first = card_num;
    card_num += 24;
    xhttp.open("GET", getBaseUrl() + "xml" + window.location.search + "&f=" + first + "&l=" + card_num, true);
    xhttp.send();
}
