console.log("javascript loaded..");
const infoDiv = document.querySelector("#info-div");
const tableDiv = document.querySelector("#table-div");
const info = document.querySelector("#info");
infoDiv.style.display="none";
tableDiv.style.display="none";

// jquery
$(document).ready(function(){

    $(".form-floating").submit(function(event) {
        event.preventDefault();
        tableDiv.style.display="none";
        infoDiv.style.display="block";
        info.innerHTML = "Searching...";
        var path = $("#floatingInputValue").val();
        if (!path) {
            return
        }
        $.post("/", {path:path})
            .done(function(xhr){
                // clear table
                $('.table > tbody:last').empty();
                // set table
                info.innerHTML = "Result";
                tableDiv.style.display="block";
                for (var i = 0; i < xhr["data"].length; i++) {
                    $('.table > tbody:last').append('<tr><td>' + String(i+1) + '</td><td>' + xhr["data"][i]["path"] + '</td><td>' + xhr["data"][i]["framein"] + '</td><td>' + xhr["data"][i]["frameout"] + '</td><td>' + xhr["data"][i]["framerange"] + '</td><td>' + xhr["data"][i]["ext"] + '</td><td>' + xhr["data"][i]["width"] + '</td><td>' + xhr["data"][i]["height"] + '</td><td>' + xhr["data"][i]["fps"] + '</td><td>' + xhr["data"][i]["codec"] + '</td></tr>')
                };
            })
            .fail(function(xhr) {
                alert(xhr.responseText)
            });
    });
});


// vanilla javascript
const pathForm = document.querySelector(".basic-form");
const pathInput = document.querySelector("#floatingInput");

function handleSubmit(event) {
    event.preventDefault();
    tableDiv.style.display="none";
    infoDiv.style.display="block";
    info.innerHTML = "Searching...";
    const path = pathInput.value;
    if (!path) {
        return
    }
    fetch("/", {
        method: "POST",
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: new URLSearchParams({
            path:path
        })
    })
    .then(response => response.json())
    .then(data => setTable(data))
    .catch(error => console.log(error)); // error catch failed....
};

function setTable(data) {
    if ("error" in data) {
        alert(data["error"])
        return
    }
    fileList = data["data"]
    let table = document.querySelector(".table");
    let tbodyRef = document.querySelector(".table").getElementsByTagName('tbody')[0];
    // clear table
    for (var i = table.rows.length - 1; i > 0; i--){
        table.deleteRow(i);
    }
    // set table
    info.innerHTML = "Result : ";
    tableDiv.style.display="block";
    for (var i = 0; i < fileList.length; i++) {
        let newRow = tbodyRef.insertRow();
        newRow.insertCell().appendChild(document.createTextNode(i+1));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["path"]));
        newRow.insertCell().appendChild(document.createTextNode(fillZero(fileList[i]["pad"], fileList[i]["framein"])));
        newRow.insertCell().appendChild(document.createTextNode(fillZero(fileList[i]["pad"], fileList[i]["frameout"])));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["framerange"]));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["ext"]));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["width"]));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["height"]));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["fps"]));
        newRow.insertCell().appendChild(document.createTextNode(fileList[i]["codec"]));
    }
};

function fillZero(p, n){
    let width = String(p);
    let num = String(n);
    return num.length >= width ? num:new Array(width-num.length+1).join('0')+num;//남는 길이만큼 0으로 채움
}

pathForm.addEventListener('submit', handleSubmit);