
const search_input = document.getElementById("search-input");
const search_image = document.getElementById("search-image");

search_input.addEventListener('keyup', InputHandler);
search_image.addEventListener('click', ImageHandler);

function InputHandler (key_event) {
    if (event.keyCode === 13) {
        location.replace("results?q="+search_input.value);
    }
}

function ImageHandler () {
    location.replace("results?q="+search_input.value);
}
