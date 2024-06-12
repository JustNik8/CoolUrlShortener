let shortUrlElem = document.getElementById("short_url");
let urlInput = document.getElementById("url_input")

let serverDomain = ""
setupServerDomain()

function setupServerDomain() {
    fetch("domain.txt")
        .then(response => response.text())
        .then(text => {
            serverDomain = text
        })

}

async function onShortenerClick() {
    let longURL = urlInput.value
    let url = `http://${serverDomain}/api/save_url`

    let response = await fetch(url, {
        method: "POST",
        body: JSON.stringify({
            long_url: longURL
        }),
    });


    let json = await response.json();

    shortUrlElem.style.visibility = "visible";
    shortUrlElem.innerHTML = json.short_url;
}


function copyToClipboard() {
    let copyText = shortUrlElem.innerHTML

    // Copy the text inside the text field
    navigator.clipboard.writeText(copyText)
        .then(() => console.log("Copied"))
        .catch(err => console.log(err));

}