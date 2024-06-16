let shortUrlElem = document.getElementById("short_url");
let urlInput = document.getElementById("url_input")

let serverDomain = ""
setupServerDomain()

function setupServerDomain() {
    fetch("domain.txt")
        .then(response => response.text())
        .then(text => {
            console.log(serverDomain)
            serverDomain = text
        })
        .then(() => loadTopUrl())

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

async function loadTopUrl() {
    let url = `http://${serverDomain}/api/top_urls?page=1&limit=10`

    let response = await fetch(url);

    let json = await response.json()

    topUrls = json.top_url_data

    let table = document.getElementById("top_urls_table")

    for (let i = 0; i < topUrls.length; i++) {
        let row = table.insertRow(i + 1);

        let sourceUrlCell = row.insertCell(0);
        let shortUrlCell = row.insertCell(1);
        let followCountCell = row.insertCell(2)
        let createCountCell = row.insertCell(3)

        let longUrl = topUrls[i].long_url
        let shortUrl = `http://${serverDomain}/${topUrls[i].short_url}`

        let longUrlElem = document.createElement('a')
        longUrlElem.setAttribute('href', longUrl)
        longUrlElem.innerHTML = longUrl

        sourceUrlCell.appendChild(longUrlElem)
        shortUrlCell.innerHTML = shortUrl
        followCountCell.innerHTML = topUrls[i].follow_count
        createCountCell.innerHTML = topUrls[i].create_count
    }
}


function copyToClipboard() {
    let copyText = shortUrlElem.innerHTML

    // Copy the text inside the text field
    navigator.clipboard.writeText(copyText)
        .then(() => console.log("Copied"))
        .catch(err => console.log(err));
}