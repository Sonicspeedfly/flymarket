async function installment() {
     var months = document.getElementById("months")
    var valuem = months.value
    var category = document.getElementById("category")
    var valuec = category.value
    var price = document.getElementById("price").value
    var score = decodeURIComponent(location.search.substr(1)).split('&');
    score.splice(0, 1);
    result = score[0]
    url = 'http://127.0.0.1:9999/api/products/installment?id=' + '1' + '&months=' + valuem + '&category=' + valuec + '&price=' + price;
    response = await fetch(url)
    content = await response.json()
    list = document.querySelector(".result")
    if (content > 0) {
        list.innerHTML += `<p classs="complete">Общая сумма рассрочки составляет : ${content}</p>`
    }
}

